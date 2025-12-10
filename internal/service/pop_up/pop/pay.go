package pop

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/model/api"
	"cngamesdk.com/api/model/cache/user"
	"context"
	"fmt"
	"go.uber.org/zap"
)

type PayPopUpService struct {
	Req api.BuildPopUpReq
}

func (receiver *PayPopUpService) Show(ctx context.Context) (config api.PopUpConfig) {
	userId := receiver.Req.UserId
	userModel := user.NewOdsUserInfoLogModel()
	userInfoErr := userModel.TakeById(ctx, receiver.Req.UserId)
	if userInfoErr != nil {
		global.Logger.ErrorCtx(ctx, "弹窗获取用户信息异常", zap.Error(userInfoErr))
		return
	}
	//实名(每次都弹)
	if userModel.TrueName == "" {
		config = api.PopUpConfig{
			Show: 1,
			Url:  "https://www.baidu.com/",
			Btns: []api.PopUpConfigBtn{
				{Type: api.PopUpConfigBtnRealName, Text: "立即实名"},
				{Type: api.PopUpConfigBtnCancel, Text: "知道了"},
			},
		}
		return
	}

	//防沉迷

	//其他弹窗
	cacheKey := fmt.Sprintf("cache-pop-pay-%d", userId)
	existsData, existsErr := global.MyRedis.Exists(ctx, cacheKey).Result()
	if existsErr != nil {
		global.Logger.ErrorCtx(ctx, "缓存异常", zap.Error(existsErr))
		return
	}
	if existsData > 0 {
		return
	}
	config = api.PopUpConfig{
		Show: 1,
		Url:  "https://www.baidu.com/",
		Btns: []api.PopUpConfigBtn{
			{Type: api.PopUpConfigBtnConfirm, Text: "知道了"},
		},
	}
	global.MyRedis.Set(ctx, cacheKey, config, PopUpRule24Hours())
	return
}
