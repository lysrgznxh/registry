package api

import (
	"net/http"
)

// 云端开放接口汇总
func localApiRouter(mux *http.ServeMux) {

	//分类列表
	mux.HandleFunc("/v1/category/list", categoryList)

	//agent 注册
	mux.HandleFunc("/v1/agent/reg", agentRegister)

	//agent图片上传
	mux.HandleFunc("/v1/agent/logo/upload", agentLogoUpload)

	//mcp 列表
	mux.HandleFunc("/v1/agent/list", agentList)

	fs := http.FileServer(http.Dir("./uploads")) //文件真实路径

	//资源映射路径
	mux.Handle("/logos/", http.StripPrefix("/logos/", fs))

}
