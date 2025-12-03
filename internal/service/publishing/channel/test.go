package channel

import (
	"cngamesdk.com/api/model/api"
	"context"
)

type TestChannel struct {
	baseChannel
}

func (receiver TestChannel) GetId() int64 {
	return 1
}

func (receiver TestChannel) Login(ctx context.Context, req *api.PublishingLoginReq) (openId string, err error) {
	return
}

func (receiver TestChannel) Pay(ctx context.Context, req *api.PublishingPayReq) (resp string, err error) {
	return
}

func (receiver TestChannel) PayCallback(ctx context.Context, req api.PublishingChannelPayCallbackReq) (
	resp api.PayCallbackData, err error) {
	return
}
