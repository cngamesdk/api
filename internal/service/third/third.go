package third

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/internal/service/user"
	"cngamesdk.com/api/model/api"
	"cngamesdk.com/api/model/sql/publishing"
	user2 "cngamesdk.com/api/model/sql/user"
	"context"
	error2 "github.com/cngamesdk/go-core/model/error"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ThirdService 第三方服务
type ThirdService struct {
}

// 第三方用户授权
func (receiver *ThirdService) Login(ctx context.Context, req *api.PublishingLoginReq) (resp *user2.OdsUserInfoLogModel, err error) {
	userService := &user.UserService{}
	bindLogModel := publishing.NewOdsPublishingUserBindLogModel()
	var userModel *user2.OdsUserInfoLogModel
	if takeErr := bindLogModel.Take(ctx, "*", "platform_id = ? and channel_id = ? and open_id = ?", req.PlatformId, req.ChannelId, req.OpenId); takeErr != nil {
		if !errors.Is(error2.ErrorRecordIsNotFind, takeErr) {
			err = takeErr
			global.Logger.ErrorCtx(ctx, "获取绑定异常", zap.Any("err", takeErr))
			return
		}
		accountAllocateReq := &api.AccountAllocateReq{}
		accountAllocateResult, accountAllocateErr := userService.AccountAllocate(ctx, accountAllocateReq)
		if accountAllocateErr != nil {
			err = accountAllocateErr
			global.Logger.ErrorCtx(ctx, "分配账号异常", zap.Any("err", accountAllocateErr))
			return
		}
		accountRegReq := &api.AccountRegReq{}
		accountRegReq.CommonReq = req.CommonReq
		accountRegReq.UserName = accountAllocateResult.UserName
		accountRegReq.Password = accountAllocateResult.Password
		regResult, regErr := userService.AccountReg(ctx, accountRegReq)
		if regErr != nil {
			err = regErr
			global.Logger.Error("注册异常", zap.Any("err", regErr))
			return
		}
		userModel = regResult
	} else { //老用户
		userModelTmp := user2.NewOdsUserInfoLogModel()
		if takeErr2 := userModelTmp.Take(ctx, "*", "id = ?", bindLogModel.UserId); takeErr2 != nil {
			err = takeErr2
			global.Logger.ErrorCtx(ctx, "获取用户异常", zap.Any("err", takeErr2))
			return
		}
		userModel = userModelTmp
	}
	resp = userModel
	return
}
