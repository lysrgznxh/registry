package api

import (
	"agent-network-protocol/registry/core"
	"agent-network-protocol/registry/util/serializer"
	"net/http"
)

// 分类列表
func categoryList(w http.ResponseWriter, r *http.Request) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmServApi)
	// 设置跨域支持
	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有来源访问
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	list := map[int]string{
		1: "基础",
		2: "营销",
		3: "数据分析",
	}
	log.Info("go")
	serializer.SuccessRes(w, list)
}
