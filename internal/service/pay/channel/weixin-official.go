package channel

import (
	"cngamesdk.com/api/model/api"
	"context"
)

type WeiXinOfficialPayChannel struct {
	WeiXinPayChannel
}

func NewWeiXinOfficialPayChannel() *WeiXinOfficialPayChannel {
	model := &WeiXinOfficialPayChannel{}
	model.WeiXinPayChannel.GetConfig = model.GetConfig
	return model
}

func (receiver *WeiXinOfficialPayChannel) GetPayChannelId() int {
	return 1
}

func (receiver *WeiXinOfficialPayChannel) GetConfig() WeiXinPayConfig {
	return WeiXinPayConfig{
		AppId:       "test",
		CallbackUrl: "https://www.baidu.com/",
	}
}

func (receiver *WeiXinOfficialPayChannel) PreOrder(ctx context.Context, req PreOrderReq) (resp PreOrderResp, err error) {
	resp, err = receiver.WeiXinPayChannel.PreOrder(ctx, req)
	return
}

func (receiver *WeiXinOfficialPayChannel) Callback(ctx context.Context, req api.PayChannelCallbackReq) (
	resp api.PayCallbackData, err error) {
	return
}
