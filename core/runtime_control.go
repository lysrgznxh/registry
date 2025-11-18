package core

import (
	"agent-network-protocol/registry/util/dir"
	"agent-network-protocol/registry/util/log"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// 运行控制公共变量,可通过包名直接使用,该变量的初始化由
var RunCtrl = NewRuntimeControl(context.Background())

func NewRuntimeControl(ctx context.Context) *RuntimeControl {
	programPath, _ := filepath.Abs(os.Args[0])

	//应用程序.exe 所在的目录
	executePath, _ := filepath.Split(programPath)

	userHome, err := dir.Home()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(ctx)
	runCtrl := &RuntimeControl{
		ThreadSize:      0,
		Wg:              &sync.WaitGroup{},
		Ctx:             ctx,
		IsCancel:        false,
		cancel:          cancel,
		RuntimePath:     executePath,
		MysqlPath:       filepath.Join(executePath, MysqlDir),
		LogPath:         filepath.Join(executePath, ClientLogDir),
		UserHome:        userHome,
		ChainSyncState:  &ChainSyncState{DbHeight: 0, HeightGap: -1},
		runtimeLogs:     make(chan ThreadLog, 40),
		TimeDiffWithNtp: 0,
	}

	runCtrl.ClientLogPath = filepath.Join(runCtrl.LogPath, "node.log")
	return runCtrl
}

type ChainSyncState struct {
	DbHeight       int64 //持久化到数据库中的高度
	DbBlockTime    int64 //持久化到数据库中的区块时间
	ChainBlockTime int64 //链上区块时间
	HeightGap      int64 //本地高度和链上高度的差距
	ChainHeight    int64 //链上的高度
	CatchUp        bool  //是否达到最新高度
}

/*
*
检查高度落差
*/
func (this *ChainSyncState) CheckHeightGap() bool {
	if this.HeightGap > 5 {
		return false
	} else if this.HeightGap == -1 {
		return false
	}
	return true
}

// 运行控制
type RuntimeControl struct {
	ThreadSize      int                //线程数量
	Wg              *sync.WaitGroup    //线程控制
	Ctx             context.Context    //主线程用来通知子线程是否需要退出
	ChainSyncState  *ChainSyncState    //主链同步状态
	IsCancel        bool               //主线程是否发出退出信号
	cancel          context.CancelFunc //释放退出信号的方法
	RuntimePath     string             // 在那个路径下启动的程序,和程序所在目录不一定是一样的
	MysqlPath       string             //mysql所在路径
	LogPath         string             //日志存放路径
	ClientLogPath   string             //网关日志路径
	UserHome        string             //用户目录,用来存储秘钥对数据库
	cpuLoads        []int              //cpu负载
	runtimeLogs     chan ThreadLog     //运行日志
	TimeDiffWithNtp int64              //与ntp服务的时差，单位毫秒
}

type ThreadLog struct {
	Name   string
	Status int
}

// 返回线程运行日志的通道
func (this *RuntimeControl) GetThreadLogChannel() chan ThreadLog {
	return this.runtimeLogs
}

// 保存启动信息  1 表示启动  0 表示退出
func (this *RuntimeControl) AddThreadLog(name string, status int) {
	log := BuildLog("", LmRunCtrl)
	if status == 1 {
		log.WithFields(logrus.Fields{"service": name, "total": this.ThreadSize}).Info("Starting")
	} else if status == 0 {
		this.ThreadSize--
		log.WithFields(logrus.Fields{"service": name, "surplus": this.ThreadSize}).Error("Exited")
	}
	this.runtimeLogs <- ThreadLog{
		Name:   name,
		Status: status,
	}
}

// 关闭主程序
// err 因为什么错误而关闭程序
func (this *RuntimeControl) Close(err interface{}) {
	this.WritePanicFile(err)
	this.Cancel()
}

// 获取前端的版本号
func (this *RuntimeControl) GetForeendVersion() string {
	file := filepath.Join(this.RuntimePath, "version")
	_, err := os.Stat(file)
	if err != nil {
		return ""
	}
	contentBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}
	return string(contentBytes)
}

// 写入崩溃日志文件
func (this *RuntimeControl) WritePanicFile(writeErr interface{}) {
	log := BuildLog("", LmRunCtrl)
	log.Debug("Save Panic Log:", writeErr)
	//检查文件是否存在的方法
	checkSourceExist := func(path string) (bool, error) {
		_, err := os.Stat(path)
		if err == nil {
			return true, nil
		}
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	//日志主目录不存在时，重建
	if ok, _ := checkSourceExist(this.LogPath); !ok {
		// 创建文件夹
		if err := os.MkdirAll(this.LogPath, os.ModePerm); err != nil {
			log.WithError(err).WithField("path", this.LogPath).Error("os.MkdirAll")
			panic(err)
		}
	}
	openFile1 := filepath.Join(this.LogPath, "panic.log")
	errFile, err := os.OpenFile(openFile1, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.WithError(err).WithField("file", openFile1).Error("os.OpenFile")
		panic(err)
		return
	}
	defer errFile.Close()
	if writeErr == nil {
		return
	}
	errFile.Write([]byte(ErrorFormat(fmt.Sprintf("%v", writeErr))))
	return
}

// 错误信息匹配
func ErrorFormat(err string) string {
	if strings.Contains(err, "Cannot allocate memory for the buffer pool") {
		return "Out of available memory"
	}
	return err
}

// 给子线程释放退出信号
func (this *RuntimeControl) Cancel() {
	log.Info("Exit notification!!!")
	this.IsCancel = true
	this.cancel()
}

// 等待需要停止的信号
func (this *RuntimeControl) NeedStop() <-chan struct{} {
	return this.Ctx.Done()
}
