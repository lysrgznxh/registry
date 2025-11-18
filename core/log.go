package core

import (
	"agent-network-protocol/registry/util/log"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"sort"
	"strings"
	"time"
)

// 日志模块使用 缩减的前缀 Lm = LogModel
var (
	LmServNtp  = log.RegisterModule("a-ntp", logrus.InfoLevel)  //时间校验服务
	LmServApi  = log.RegisterModule("a-api", logrus.InfoLevel)  //api
	LmRunCtrl  = log.RegisterModule("a-rct", logrus.InfoLevel)  //运行控制
	LmMain     = log.RegisterModule("a-m", logrus.InfoLevel)    //入口
	LmNodeUtil = log.RegisterModule("a-u", logrus.InfoLevel)    //节点 util包
	LmCore     = log.RegisterModule("a-core", logrus.InfoLevel) //节点 db/chain 包
	LmDbUser   = log.RegisterModule("a-dbu", logrus.InfoLevel)  //节点 db/user 包
)

// 构建日志
var BuildLog = log.BuildLog

// 解析环境变量
func InitLogger(logPath string, enableLogSave bool, defaultLevel logrus.Level) {
	// 初始化日志log,默认info
	log.InitLogger(defaultLevel)

	if enableLogSave {
		logSaveTime := time.Hour * 24 * 7
		logSplitTime := time.Hour * 24
		log.EnableLogStorage(logPath, logSaveTime, logSplitTime)
	}

	//默认日志级别,对所有模块生效
	//范例: export DST_LOGGING="debug"
	logDefault, ok := os.LookupEnv("DST_LOGGING")
	if ok {
		logDefault = strings.ReplaceAll(logDefault, "\"", "")
		level, err := logrus.ParseLevel(logDefault)
		fmt.Println("env define default-log-level:", logDefault)
		if err == nil { //如果没有错误
			log.ResetAllModuleLevel(level) //重置所有模块的日志级别
			fmt.Println("used default-log-level:", level)
		}
	}

	//检查是否设置了对象模块的日志级别
	//范例:   默认级别为info,指定的模块级别为debug： export DST_LOGGING="info;module1:debug;module2:debug;"
	chainLogging, ok := os.LookupEnv("DST_LOGGING")
	if ok {
		chainLogging = strings.Trim(chainLogging, `"`)
		sets := strings.Split(chainLogging, ";") //根据;组装成数组
		for _, set := range sets {               //格式 xx1=xxx xx2=xxx
			//fmt.Println("set=", set)
			kvs := strings.Split(set, ":")
			//fmt.Println("length;", len(kvs))
			if len(kvs) != 2 {
				continue
			}
			moduleStr := kvs[0]
			levelStr := kvs[1]

			logModule := log.LogModule(moduleStr)
			//如果指定的日志模块并不存在,忽略
			if !log.CheckModuleExists(logModule) {
				continue
			}
			level, err := logrus.ParseLevel(levelStr)
			if err == nil { //如果没有错误
				log.SetModuleLevel(logModule, level) //更新模块日志级别
			}
		}
		fmt.Println("module-log-level--->")
		var keys []string
		//先排序
		models := log.GetModelLevels()
		for k, _ := range models {
			keys = append(keys, string(k))
		}
		sort.Strings(keys)

		//然后输出
		for _, moduleName := range keys {
			fmt.Println(moduleName, ":", models[log.LogModule(moduleName)])
		}
		fmt.Println("<---")
	}
}
