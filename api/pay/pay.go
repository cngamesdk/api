package pay

import (
	"cngamesdk.com/api/internal/logic/pay"
	"cngamesdk.com/api/model/api"
	response2 "cngamesdk.com/api/model/api/response"
	"github.com/cngamesdk/go-core/model/response"
	"github.com/cngamesdk/go-core/translator"
	"github.com/gin-gonic/gin"
	"net/http"
)

// PayChannelCallback 支付渠道回调
func PayChannelCallback(ctx *gin.Context) {
	var req api.PayChannelCallbackReq
	if err := ctx.ShouldBind(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	resp, err := (&pay.PayLogic{}).PayChannelCallback(ctx, req)
	msg := resp
	if err != nil {
		msg = err.Error()
	}
	ctx.String(http.StatusOK, msg)
	return
}

// PublishingChannelPayCallback 发行渠道支付回调
func PublishingChannelPayCallback(ctx *gin.Context) {
	var req api.PublishingChannelPayCallbackReq
	if err := ctx.ShouldBind(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	resp, err := (&pay.PayLogic{}).PublishingChannelPayCallback(ctx, req)
	msg := resp
	if err != nil {
		msg = err.Error()
	}
	ctx.String(http.StatusOK, msg)
	return
}
