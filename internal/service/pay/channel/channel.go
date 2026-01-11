package channel

import (
	"cngamesdk.com/api/model/api"
	"context"
	"github.com/cngamesdk/go-core/pay"
	"github.com/duke-git/lancet/v2/random"
	"time"
)

type PreOrderReq struct {
	Money int
}

type PreOrderResp struct {
	OrderId string
	Url     string
}

type WeiXinPayConfig struct {
	AppId       string
	CallbackUrl string
	// other
}

type AlipayPayConfig struct {
	AppId       string
	CallbackUrl string
	// other
}

type PayChannelInterface interface {
	GetPayChannelId() int64
	PreOrder(ctx context.Context, req PreOrderReq) (resp PreOrderResp, err error)
	Callback(ctx context.Context, req api.PayChannelCallbackReq) (resp api.PayCallbackData, err error)
}

type PayChannel struct {
}

func (receiver PayChannel) CreateOrderId() string {
	return time.Now().Format("20060102150405") + random.RandNumeral(10)
}

type WeiXinPayChannel struct {
	PayChannel
	GetConfig func() WeiXinPayConfig
	PayChannelInterface
}

func (receiver WeiXinPayChannel) PreOrder(ctx context.Context, req PreOrderReq) (resp PreOrderResp, err error) {
	config := receiver.GetConfig()
	orderId := receiver.CreateOrderId()
	payService := &pay.WeiXinPay{}
	preOrderReq := pay.PreOrderReq{}
	preOrderReq.Money = req.Money
	preOrderReq.OrderId = orderId
	preOrderReq.CallbackUrl = config.CallbackUrl
	preOrderResult, preOrderErr := payService.PreOrder(ctx, preOrderReq)
	if preOrderErr != nil {
		err = preOrderErr
		return
	}
	resp.OrderId = orderId
	resp.Url = preOrderResult.Url
	return
}

type AlipayPayChannel struct {
	PayChannel
	GetConfig func() AlipayPayConfig
	PayChannelInterface
}

func (receiver AlipayPayChannel) PreOrder(ctx context.Context, req PreOrderReq) (resp PreOrderResp, err error) {
	config := receiver.GetConfig()
	orderId := receiver.CreateOrderId()
	payService := &pay.Alipay{}
	preOrderReq := pay.PreOrderReq{}
	preOrderReq.Money = req.Money
	preOrderReq.OrderId = orderId
	preOrderReq.CallbackUrl = config.CallbackUrl
	preOrderResult, preOrderErr := payService.PreOrder(ctx, preOrderReq)
	if preOrderErr != nil {
		err = preOrderErr
		return
	}
	resp.OrderId = orderId
	resp.Url = preOrderResult.Url
	return
}
