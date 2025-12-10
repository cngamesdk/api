package pop_up

import (
	"cngamesdk.com/api/internal/service/pop_up/pop"
	"cngamesdk.com/api/model/api"
	"context"
)

func GetPopUpFactory() {

}

// PopUpService 弹窗服务
type PopUpService struct {
}

// GetPopUpConfig 获取弹窗配置
func (receiver *PopUpService) GetPopUpConfig(ctx context.Context, req api.BuildPopUpReq) (resp api.PopUpConfig) {
	var popServiceInterface pop.PopUpInterface
	switch req.Source {
	case api.BuildPopUpSourceInit:
		popService := &pop.InitPopUpService{}
		popService.Req = req
		popServiceInterface = popService
		break
	case api.BuildPopUpSourceLogin:
		popService := &pop.LoginPopUpService{}
		popService.Req = req
		popServiceInterface = popService
		break
	case api.BuildPopUpSourcePay:
		popService := &pop.PayPopUpService{}
		popService.Req = req
		popServiceInterface = popService
		break
	case api.BuildPopUpSourceHeartbeat:
		popService := &pop.HeartbeatPopUpService{}
		popService.Req = req
		popServiceInterface = popService
		break
	default:
		break
	}
	if popServiceInterface != nil {
		resp = popServiceInterface.Show(ctx)
	}
	return
}
