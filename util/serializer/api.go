package serializer

type ApiCommon struct {
	Status int    `json:"status"`
	Info   string `json:"info"`
}

type ApiPageCommon struct {
	Status int    `json:"status"`
	Info   string `json:"info"`
	Total  int    `json:"total"`
}

type ApiReturnString struct {
	Status string `json:"status"`
	Info   string `json:"info"`
	Data   string `json:"data"`
}

type ApiReturn struct {
	ApiCommon
	Data string `json:"data"`
}

type ApiPageReturn struct {
	ApiPageCommon
	Data interface{} `json:"data"`
}

type ApiReturnInterface struct {
	ApiCommon
	Data interface{} `json:"data"`
}

type ApiReturnArray struct {
	ApiCommon
	Data []string `json:"data"`
}
