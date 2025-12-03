package api

import (
	"cngamesdk.com/api/internal/logic/api"
	api2 "cngamesdk.com/api/model/api"
	response2 "cngamesdk.com/api/model/api/response"
	"github.com/cngamesdk/go-core/model/code"
	"github.com/cngamesdk/go-core/model/response"
	"github.com/cngamesdk/go-core/translator"
	"github.com/gin-gonic/gin"
)

// PublishingLogin 发行渠道登录
func PublishingLogin(ctx *gin.Context) {
	var req api2.PublishingLoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api.PublishingLogic{}).Login(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()).SetData(resp))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// PublishingPay 发行渠道支付
func PublishingPay(ctx *gin.Context) {
	var req api2.PublishingPayReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api.PublishingLogic{}).Pay(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()).SetData(resp))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}
