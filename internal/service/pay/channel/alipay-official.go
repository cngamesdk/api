package channel

import (
	"cngamesdk.com/api/model/api"
	"context"
)

type AlipayOfficialPayChannel struct {
	AlipayPayChannel
}

func NewAlipayOfficialPayChannel() *AlipayOfficialPayChannel {
	model := &AlipayOfficialPayChannel{}
	model.AlipayPayChannel.GetConfig = model.GetConfig
	return model
}

func (receiver *AlipayOfficialPayChannel) GetPayChannelId() int64 {
	return 2
}

func (receiver *AlipayOfficialPayChannel) GetConfig() AlipayPayConfig {
	return AlipayPayConfig{
		AppId:       "test",
		CallbackUrl: "https://www.xxx.com/",
	}
}

func (receiver *AlipayOfficialPayChannel) PreOrder(ctx context.Context, req PreOrderReq) (
	resp PreOrderResp, err error) {
	resp, err = receiver.PreOrder(ctx, req)
	return
}

func (receiver *AlipayOfficialPayChannel) Callback(ctx context.Context, req api.PayChannelCallbackReq) (
	resp api.PayCallbackData, err error) {
	return
}
