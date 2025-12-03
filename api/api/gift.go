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

// GiftList 礼包列表
func GiftList(ctx *gin.Context) {
	var req api2.GiftListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api.GiftLogic{}).List(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// GiftClaim 礼包领取
func GiftClaim(ctx *gin.Context) {
	var req api2.GiftClaimReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api.GiftLogic{}).Claim(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}
