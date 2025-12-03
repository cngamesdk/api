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

// LaunchDataReport 启动数据上报
func LaunchDataReport(ctx *gin.Context) {
	var req api2.LaunchReportReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api.DataReportLogic{}).Launch(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// GameLogDataReport 游戏日志数据上报
func GameLogDataReport(ctx *gin.Context) {
	var req api2.GameReportReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api.DataReportLogic{}).GameLog(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}
