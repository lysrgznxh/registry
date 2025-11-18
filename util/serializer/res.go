package serializer

import (
	"encoding/json"
	"net/http"
)

// 直接返回response不进行格式包装
func DirectResponse(w http.ResponseWriter, data interface{}) {
	resByte, _ := json.Marshal(data)
	headerSet(w)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resByte)
	return
}

func headerSet(w http.ResponseWriter) {
	Cors(w) //设置跨域
	w.Header().Set("Content-Type", "application/json")
}
func SuccessRes(w http.ResponseWriter, data interface{}) {
	res := ApiReturnInterface{
		ApiCommon: ApiCommon{
			Status: 1,
			Info:   "",
		},
		Data: data,
	}
	if data == nil {
		res.Data = struct{}{}
	}
	resByte, _ := json.Marshal(res)
	headerSet(w)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resByte)
	return
}

func SuccessResCall(w http.ResponseWriter, data interface{}, call func(w http.ResponseWriter)) {
	res := ApiReturnInterface{
		ApiCommon: ApiCommon{
			Status: 1,
			Info:   "",
		},
		Data: data,
	}
	if data == nil {
		res.Data = struct{}{}
	}
	resByte, _ := json.Marshal(res)
	if call != nil {
		call(w)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resByte)
	return
}
func SuccessInfoRes(w http.ResponseWriter, data interface{}, msg string) {
	res := ApiReturnInterface{
		ApiCommon: ApiCommon{
			Status: 1,
			Info:   msg,
		},
		Data: data,
	}
	if data == nil {
		res.Data = struct{}{}
	}
	resByte, _ := json.Marshal(res)
	headerSet(w)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resByte)
	return
}

func SuccessTotalRes(w http.ResponseWriter, data interface{}, total int) {
	res := ApiPageReturn{
		ApiPageCommon: ApiPageCommon{
			Status: 1,
			Total:  total,
			Info:   "",
		},
		Data: data,
	}
	if data == nil {
		res.Data = struct{}{}
	}
	resByte, _ := json.Marshal(res)
	headerSet(w)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resByte)
	return
}

func ErrorRes(w http.ResponseWriter, err error) {
	res := ApiReturnInterface{
		ApiCommon: ApiCommon{
			Status: 0,
			Info:   err.Error(),
		},
		Data: struct{}{},
	}
	resByte, _ := json.Marshal(res)
	headerSet(w)
	_, _ = w.Write(resByte)
	return
}

func StringErrorRes(w http.ResponseWriter, err error) {
	res := ApiReturnInterface{
		ApiCommon: ApiCommon{
			Status: 0,
			Info:   err.Error(),
		},
		Data: struct {
		}{},
	}
	resByte, _ := json.Marshal(res)
	headerSet(w)
	_, _ = w.Write(resByte)
	return
}
