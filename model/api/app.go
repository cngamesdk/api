package api

// AppReq 手游的全局参数
type AppReq struct {
	Data string `json:"data" form:"data" binding:"required"`
	Key  string `json:"key" form:"key" binding:"required"`
}

type AppResp struct {
	Data      string `json:"data"`
	Key       string `json:"key"`
	RequestId string `json:"request_id"`
}
