package agent

import (
	"agent-network-protocol/registry/core"
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"time"
)

func checkMysql() error {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmDbUser)
	db, err := sql.Open("mysql", core.GetMysqlDbPasswordConnInfo())
	if err != nil {
		log.WithError(err).Error("sql.Open")
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	err = db.PingContext(ctx) //这里使用带超时的上下文,防止阻塞在这里
	if err != nil {
		if mysqlError, ok := err.(*mysql.MySQLError); ok {
			if mysqlError.Number == 1045 {
				db, err = sql.Open("mysql", core.GetMysqlDbConnInfo())
				if err != nil {
					log.WithError(err).Error("sql.Open")
					return err
				}
				err = db.PingContext(ctx) //这里使用带超时的上下文,防止阻塞在这里
				if err != nil {
					log.WithError(err).Error("PingContext")
					return err
				}
			}
		}
	}
	db.Close() //关闭测试连接
	return nil
}
