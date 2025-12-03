package api

import (
	"cngamesdk.com/api/internal/logic/api"
	media2 "cngamesdk.com/api/model/api/media"
	response2 "cngamesdk.com/api/model/api/response"
	"github.com/cngamesdk/go-core/model/code"
	"github.com/cngamesdk/go-core/model/response"
	"github.com/cngamesdk/go-core/translator"
	"github.com/gin-gonic/gin"
)

// AdvertisingClick 广告点击
func AdvertisingClick(ctx *gin.Context) {

}

// ReportReg 上报注册
func ReportReg(ctx *gin.Context) {
	var req media2.ReportRegReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api.ReportLogic{}).Reg(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// ReportLogin 上报登录
func ReportLogin(ctx *gin.Context) {
	var req media2.ReportLoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api.ReportLogic{}).Login(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// ReportPay 上报付费
func ReportPay(ctx *gin.Context) {
	var req media2.ReportPayReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api.ReportLogic{}).Pay(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// ReportCallback 上报回调
func ReportCallback(ctx *gin.Context) {
	var req media2.ReportCallbackReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api.ReportLogic{}).Callback(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}
