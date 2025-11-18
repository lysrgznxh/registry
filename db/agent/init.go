package agent

import (
	"agent-network-protocol/registry/core"
	"github.com/go-xorm/xorm"
)

var agentDb *xorm.Engine

func GetAgentDb() *xorm.Engine {
	return agentDb
}

// 初始化agent数据库
func Init(runCtrl *core.RuntimeControl) error {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmDbUser)
	var err error
	err = checkMysql()
	if err != nil {
		log.WithError(err).Error("checkMysql")
		return err
	}

	agentDb, err = xorm.NewEngine("mysql", core.GetMysqlDbPasswordConnInfo())
	if err != nil {
		log.WithError(err).Error("NewEngine")
		return err
	}

	err = agentDb.Sync2(new(Agent))
	if err != nil {
		log.WithError(err).Error("Coin Sync2")
		return err
	}

	log.Info("数据库加载完成")
	return err
}
