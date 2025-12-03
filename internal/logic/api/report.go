package api

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/internal/service/media"
	media2 "cngamesdk.com/api/model/api/media"
	"context"
	"go.uber.org/zap"
)

type ReportLogic struct {
}

func (r *ReportLogic) Reg(ctx context.Context, req *media2.ReportRegReq) (resp interface{}, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		global.Logger.ErrorCtx(ctx, "验证失败", zap.Error(validateErr))
		return
	}
	service := &media.ReportService{}
	resp, err = service.Reg(ctx, req)
	return
}

func (r *ReportLogic) Login(ctx context.Context, req *media2.ReportLoginReq) (resp interface{}, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		global.Logger.ErrorCtx(ctx, "验证失败", zap.Error(validateErr))
		return
	}
	service := &media.ReportService{}
	resp, err = service.Login(ctx, req)
	return
}

func (r *ReportLogic) Pay(ctx context.Context, req *media2.ReportPayReq) (resp interface{}, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		global.Logger.ErrorCtx(ctx, "验证失败", zap.Error(validateErr))
		return
	}
	service := &media.ReportService{}
	resp, err = service.Pay(ctx, req)
	return
}

func (r *ReportLogic) Callback(ctx context.Context, req *media2.ReportCallbackReq) (resp media2.ReportCallbackResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		global.Logger.ErrorCtx(ctx, "验证失败", zap.Error(validateErr))
		return
	}
	service := &media.ReportService{}
	respService, errService := service.Callback(ctx, req)
	if errService != nil {
		err = errService
		global.Logger.ErrorCtx(ctx, "回传异常", zap.Error(errService))
		return
	}
	resp = respService
	return
}
