package serializer

import (
	"net/http"
)

// 设置跨域设置,注意 头key定义不要重复，不然的话在部分设备下会导致跨域策略失败
func Cors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	//w.Header().Set("content-type", "application/json;charset=UTF-8")
}
