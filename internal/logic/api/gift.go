package api

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/internal/service/gift"
	"cngamesdk.com/api/model/api"
	"context"
	"go.uber.org/zap"
)

type GiftLogic struct {
}

// List 礼包列表
func (receiver *GiftLogic) List(ctx context.Context, req *api.GiftListReq) (resp api.GiftListResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	service := &gift.GiftService{}
	serviceResp, serviceErr := service.List(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("服务调用异常", zap.Any("err", serviceErr))
		return
	}
	resp = serviceResp
	return
}

// Claim 礼包领取
func (receiver *GiftLogic) Claim(ctx context.Context, req *api.GiftClaimReq) (resp api.GiftClaimResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	service := &gift.GiftService{}
	serviceResp, serviceErr := service.Claim(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("服务调用异常", zap.Any("err", serviceErr))
		return
	}
	resp = serviceResp
	return
}
