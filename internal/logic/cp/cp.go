package cp

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/internal/service/cp"
	"cngamesdk.com/api/model/api"
	"context"
	"go.uber.org/zap"
)

type CpLogic struct {
}

func (receiver CpLogic) Auth(ctx context.Context, req *api.CpAuthReq) (resp api.CpAuthResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	cpService := cp.CpService{}
	authResp, authErr := cpService.Auth(ctx, req)
	if authErr != nil {
		err = authErr
		global.Logger.ErrorCtx(ctx, "授权异常", zap.Any("err", authErr))
		return
	}
	resp = authResp
	return
}
