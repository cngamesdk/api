package api

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/internal/service/data"
	"cngamesdk.com/api/model/api"
	"context"
	"go.uber.org/zap"
)

type DataReportLogic struct {
}

// Launch 启动数据上报
func (receiver *DataReportLogic) Launch(ctx context.Context, req *api.LaunchReportReq) (
	resp api.LaunchReportResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	service := data.DataReportService{}
	serviceResp, serviceErr := service.Launch(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("服务调用异常", zap.Any("err", serviceErr))
		return
	}
	resp = serviceResp
	return
}

// GameLog 游戏日志上报
func (receiver *DataReportLogic) GameLog(ctx context.Context, req *api.GameReportReq) (
	resp api.GameReportResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	service := data.DataReportService{}
	serviceResp, serviceErr := service.Game(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("服务调用异常", zap.Any("err", serviceErr))
		return
	}
	resp = serviceResp
	return
}
