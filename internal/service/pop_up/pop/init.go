package pop

import (
	"cngamesdk.com/api/model/api"
	"context"
)

type InitPopUpService struct {
	BasePopUpService
}

func (receiver *InitPopUpService) Show(ctx context.Context) (config api.PopUpConfig) {
	config = api.PopUpConfig{
		Show: 1,
		Url:  "https://www.baidu.com/",
		Btns: []api.PopUpConfigBtn{
			{Type: api.PopUpConfigBtnConfirm, Text: "知道了"},
		},
	}
	return
}
