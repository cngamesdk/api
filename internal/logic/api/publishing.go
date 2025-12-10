package api

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/internal/service/pop_up"
	"cngamesdk.com/api/internal/service/publishing"
	"cngamesdk.com/api/internal/service/user"
	"cngamesdk.com/api/model/api"
	"context"
	"go.uber.org/zap"
)

type PublishingLogic struct {
}

// Login 发行渠道登录
func (receiver *PublishingLogic) Login(ctx context.Context, req *api.PublishingLoginReq) (
	resp api.PublishingLoginResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	service := &publishing.PublishingService{}
	serviceResp, serviceErr := service.Login(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("服务调用异常", zap.Error(serviceErr))
		return
	}

	popUpConfig := (&pop_up.PopUpService{}).GetPopUpConfig(ctx, api.BuildPopUpReq{CommonReq: req.CommonReq, UserId: serviceResp.Id, Source: api.BuildPopUpSourceLogin})
	resp.PopUp = popUpConfig

	//保存登录日志
	loginLogReq := &api.LoginLogReq{}
	loginLogReq.UserId = serviceResp.Id
	loginLogReq.CommonReq = req.CommonReq
	if saveErr := user.SaveLoginLogAsync(ctx, loginLogReq); saveErr != nil {
		global.Logger.ErrorCtx(ctx, "保存登录日志异常", zap.Any("err", saveErr))
	}

	auth, authErr := user.BuildUserAuthResp(ctx, serviceResp)
	if authErr != nil {
		err = authErr
		global.Logger.ErrorCtx(ctx, "授权异常", zap.Any("err", authErr))
		return
	}
	resp.BaseUserAuthRespModel = auth
	return
}

// Pay 发行渠道支付
func (receiver *PublishingLogic) Pay(ctx context.Context, req *api.PublishingPayReq) (
	resp api.PublishingPayResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}

	popUpConfig := (&pop_up.PopUpService{}).GetPopUpConfig(ctx, api.BuildPopUpReq{CommonReq: req.CommonReq, UserId: req.UserId, Source: api.BuildPopUpSourcePay})
	resp.PopUp = popUpConfig

	service := &publishing.PublishingService{}
	serviceResp, serviceErr := service.Pay(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("服务调用异常", zap.Error(serviceErr))
		return
	}
	resp = serviceResp
	return
}
