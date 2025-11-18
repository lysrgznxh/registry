package log

import (
	"errors"
	"github.com/sirupsen/logrus"
)

type LogModule string

//保存模块的日志对象
var logModules = make(map[string]*logrus.Logger)

//保存模块日志级别
var moduleLogLevels = make(map[LogModule]logrus.Level)

//默认日志级别 所有模块都生效
var defaultLogLevel = logrus.InfoLevel

//默认的日志对象
var Log *logrus.Logger = func() *logrus.Logger {
	log1 := logrus.New()
	log1.SetFormatter(diyTextFormatter()) //设置日志格式
	return log1
}()

//注册日志模块,如果日志模块已经存在则会panic
func RegisterModule(name string, level logrus.Level) LogModule {
	if CheckModuleExists(LogModule(name)) {
		panic("log module " + name + " only exists")
	}
	moduleLog := logrus.New() //复制log属性
	moduleLog.SetLevel(level) //设置日志级别
	moduleLog.SetFormatter(diyTextFormatter())
	logModules[name] = moduleLog
	moduleLogLevels[LogModule(name)] = level
	return LogModule(name)
}

//重置所有模块的日志级别
func ResetAllModuleLevel(level logrus.Level) (err error) {
	defaultLogLevel = level
	for moduleName, _ := range moduleLogLevels {
		err = SetModuleLevel(moduleName, level)
		if err != nil {
			return err
		}
	}
	return err
}

//设置模块日志级别
func SetModuleLevel(module LogModule, level logrus.Level) error {
	if !CheckModuleExists(module) {
		return errors.New("log module " + string(module) + " not have registerd")
	}
	logModules[string(module)].SetLevel(level)
	moduleLogLevels[module] = level
	return nil
}

//检查模块日志是否存在
func CheckModuleExists(module LogModule) bool {
	if _, ok := logModules[string(module)]; ok {
		return true
	} else {
		return false
	}
}

//返回默认日志级别
func GetDefaultLogLevel() logrus.Level {
	return defaultLogLevel
}

//获取模块的日志级别
func GetModelLevels() map[LogModule]logrus.Level {
	return moduleLogLevels
}
