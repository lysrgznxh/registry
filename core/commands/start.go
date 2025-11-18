package commands

import (
	"agent-network-protocol/registry/core"
	"agent-network-protocol/registry/db/agent"
	"agent-network-protocol/registry/services/api"
	"agent-network-protocol/registry/util/files"
	"agent-network-protocol/registry/util/net"
	"agent-network-protocol/registry/util/thread"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"
)

var (
	// 客户端日志
	logLevel = logrus.InfoLevel
)

func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start agent registry",
		Run: func(cmd *cobra.Command, args []string) {
			configFile := cmd.Flag("config").Value.String()
			StartNode(configFile)
		},
	}
	cmd.Flags().String("config", "", "config file (default current dir/config.yaml)")
	return cmd
}

// 启动节点
func StartNode(configFile string) {
	gin.SetMode(gin.ReleaseMode) //关闭gin框架的调试模式
	runCtrl := core.RunCtrl
	//检查环境变量中是否设置了日志级别，如果设置了则启用环境变量的日志级别
	core.InitLogger(runCtrl.ClientLogPath, true, logLevel)

	log := core.BuildLog(core.GetPackageFuncName(), core.LmMain)

	//该defer无法捕获线程Thread产生的错误.只能捕获非线程的panic
	defer func() {

		if err := recover(); err != nil {
			log.Error(err)
			runCtrl.WritePanicFile(err) //写入错误内容
			return
		}
	}()

	// 清空新建panic文件
	runCtrl.WritePanicFile(nil)
	

	//加载配置文件
	err := core.LoadConfig(configFile)
	if err != nil {
		log.WithError(err).Error("loadConfig")
		panic(err.Error())
	}

	fmt.Println("LocalServerUrl:", core.BaseConf.LocalServerUrl)

	//启动数据库
	err = agent.Init(runCtrl)
	if err != nil {
		log.WithError(err).Error(("init mysql failed"))
		panic(err)
	}

	fmt.Printf("\x1b[32mApiServer\x1b[0m Port:%d\n", core.APiServicePort)

	//远程接口服务
	thread.NewThreadV2("[Remote Api]", runCtrl, func(runCtrl *core.RuntimeControl) {
		api.Start(runCtrl, core.APiServicePort)
	})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//检测系统退出信号,通知子线程退出
	go func() {
		sig := <-sigs //等待退出信号
		log.Info("Exit signal received:", sig.String())
		runCtrl.Cancel() //通知子线程退出
	}()

	log.Info("Main program startup completed.")
	runCtrl.Wg.Wait()

	//延迟1秒退出,为退出消息发送留出足够的时间
	select {
	case <-time.After(time.Second * 1):
		break
	}
	log.Info("The main program has exited.")
}
