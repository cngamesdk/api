package pop

import (
	"cngamesdk.com/api/model/api"
	"context"
	"time"
)

// PopUpRule24Hours 24小时内弹窗
func PopUpRule24Hours() time.Duration {
	return time.Hour * 24
}

type BasePopUpService struct {
	Req api.BuildPopUpReq
}

type PopUpInterface interface {
	Show(ctx context.Context) (config api.PopUpConfig)
}
