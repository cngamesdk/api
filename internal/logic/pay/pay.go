package pay

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/internal/service/cp"
	"cngamesdk.com/api/internal/service/pay"
	"cngamesdk.com/api/internal/service/pay/channel"
	"cngamesdk.com/api/internal/service/publishing"
	"cngamesdk.com/api/model/api"
	"cngamesdk.com/api/model/sql/log"
	publishing2 "cngamesdk.com/api/model/sql/publishing"
	"context"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type PayLogic struct {
}

// 支付网关支付回调
func (receiver *PayLogic) PayChannelCallback(ctx context.Context, req api.PayChannelCallbackReq) (resp string, err error) {
	if validateErr := req.Validate(); validateErr != nil {
		err = validateErr
		return
	}
	var payChannel channel.PayChannelInterface
	payChannels := pay.GetPayChannels()
	for _, item := range payChannels {
		if item.GetPayChannelId() == req.GetPayChannelId() {
			payChannel = item
			break
		}
	}
	if payChannel == nil {
		err = errors.New("未找到支付渠道")
		global.Logger.Error("未找到支付渠道", zap.Any("data", req))
		return
	}
	callbackResp, callbackErr := payChannel.Callback(ctx, req)
	if callbackErr != nil {
		err = callbackErr
		global.Logger.Error("支付渠道回调异常", zap.Any("err", callbackErr))
		return
	}
	if payCallbackErr := pay.PayCallback(ctx, callbackResp); payCallbackErr != nil {
		err = payCallbackErr
		global.Logger.Error("支付回调处理异常", zap.Any("err", payCallbackErr))
		return
	}
	cpService := cp.CpService{}
	if cpErr := cpService.PayCallback(ctx, &callbackResp); cpErr != nil {
		err = cpErr
		global.Logger.Error("CP发货异常", zap.Any("err", cpErr))
		return
	}
	return
}

// PublishingChannelPayCallback 发行渠道支付回调
func (receiver *PayLogic) PublishingChannelPayCallback(ctx context.Context, req api.PublishingChannelPayCallbackReq) (resp string, err error) {
	if validateErr := req.Validate(); validateErr != nil {
		err = validateErr
		return
	}
	extData, extOk := req["ext"]
	if !extOk {
		err = errors.New("透传参数不存在")
		global.Logger.Error("透传参数不存在", zap.Any("data", req))
		return
	}
	parsePayExtResult, parsePayExtErr := publishing.ParsePayExt(cast.ToString(extData))
	if parsePayExtErr != nil {
		err = errors.Wrap(parsePayExtErr, "解析EXT异常")
		global.Logger.Error("解析EXT异常", zap.Error(parsePayExtErr))
		return
	}
	publishingChannel := publishing.GetChannel(req.GetChannelId())
	if publishingChannel == nil {
		err = errors.New("获取发行渠道失败")
		global.Logger.Error("获取发行渠道失败", zap.Any("data", req))
		return
	}
	//发行渠道支付校验
	callbackResp, callbackErr := publishingChannel.PayCallback(ctx, req)
	if callbackErr != nil {
		err = callbackErr
		global.Logger.Error("发行渠道校验失败", zap.Error(callbackErr))
		return
	}
	if callbackResp.OpenId == "" {
		err = errors.New("openid为空")
		global.Logger.Error("openid为空", zap.Any("data", callbackResp))
		return
	}

	orderModel := &log.OdsPayLogModel{}
	if takeErr := orderModel.Take(ctx, "*", "id = ?", parsePayExtResult.Id); takeErr != nil {
		err = takeErr
		global.Logger.Error("订单异常", zap.Error(takeErr), zap.Any("data", parsePayExtResult))
		return
	}
	callbackResp.PlatformId = orderModel.PlatformId
	callbackResp.OrderId = orderModel.OrderId
	bindModel := &publishing2.OdsPublishingUserBindLogModel{}
	if takeErr := bindModel.Take(ctx, "*", "platform_id = ? and channel_id = ? and open_id = ?", orderModel.PlatformId, req.GetChannelId(), callbackResp.OpenId); takeErr != nil {
		err = takeErr
		global.Logger.Error("获取异常", zap.Error(takeErr))
		return
	}

	if bindModel.UserId != orderModel.UserId {
		err = errors.New("请求非法")
		global.Logger.Error("用户不一致", zap.Any("data", orderModel), zap.Any("bind", bindModel))
		return
	}
	//支付渠道相关数据处理
	if payCallbackErr := pay.PayCallback(ctx, callbackResp); payCallbackErr != nil {
		err = payCallbackErr
		global.Logger.Error("回调异常", zap.Any("err", payCallbackErr))
		return
	}

	//发货相关处理
	cpService := cp.CpService{}
	if cpCallErr := cpService.PayCallback(ctx, &callbackResp); cpCallErr != nil {
		err = cpCallErr
		global.Logger.Error("研发异常", zap.Any("err", cpCallErr))
		return
	}
	return
}
