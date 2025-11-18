package api

import (
	"agent-network-protocol/registry/core"
	"agent-network-protocol/registry/util/net"
	"fmt"
	"net/http"
	"time"
)

/*
*
远程接口服务
*/
func Start(runCtrl *core.RuntimeControl, port int) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmServApi)

	log.Info("run")

	mux := http.NewServeMux()
	localApiRouter(mux) //开放接口路由
	errChan := make(chan error)
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), mux)
		if err != nil {
			log.Error(err.Error())
			errChan <- err //启动出错时，抛出错误给上级线程
		}
	}()

	for {
		isOk := net.TestConnectivity(fmt.Sprintf("http://127.0.0.1:%d", port), time.Second*2)
		if isOk {
			log.Info("Api Server Start Success")
			break
		}
	}

	select {
	case <-runCtrl.NeedStop(): //主线程停止时退出
		break
	case err := <-errChan:
		panic(err)
	}
}
