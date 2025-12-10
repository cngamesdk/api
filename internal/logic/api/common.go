package api

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/internal/service/pop_up"
	"cngamesdk.com/api/internal/service/user"
	"cngamesdk.com/api/model/api"
	"context"
	"github.com/cngamesdk/go-core/model/sql/common"
	errors2 "github.com/pkg/errors"
	"go.uber.org/zap"
)

type CommonLogic struct {
}

// Init 初始化
func (receiver *CommonLogic) Init(ctx context.Context, req *api.InitReq) (resp api.InitResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	resp.QqNumber = "123456"
	resp.QqLink = "mqq://im/chat?chat_type=wpa&uin=123456&version=1&src_type=web"
	resp.DomainList = []string{"https://www.xxx.com", "https://www.xxx.com"}
	resp.PayConfig = []string{
		common.PayTypeWeiXinPay,
		common.PayTypeAlipay,
	}
	resp.SwitchRealName = 1
	resp.Heartbeat = api.HeartbeatSwitch{Toggle: 1, Interval: 30}
	resp.PopUp = (&pop_up.PopUpService{}).GetPopUpConfig(ctx, api.BuildPopUpReq{CommonReq: req.CommonReq, Source: api.BuildPopUpSourceInit})
	return
}

// AccountAllocate 账号分配
func (receiver *CommonLogic) AccountAllocate(ctx context.Context, req *api.AccountAllocateReq) (resp api.AccountAllocateResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	userService := &user.UserService{}
	serviceResp, serviceErr := userService.AccountAllocate(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("分配账号异常", zap.Any("err", serviceErr))
		return
	}
	resp = serviceResp
	return
}

// AccountReg 账号注册
func (receiver *CommonLogic) AccountReg(ctx context.Context, req *api.AccountRegReq) (resp api.AccountRegResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	userService := &user.UserService{}
	serviceResp, serviceErr := userService.AccountReg(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("账号注册异常", zap.Any("err", serviceErr))
		return
	}

	loginLogReq := &api.LoginLogReq{}
	loginLogReq.UserId = serviceResp.Id
	loginLogReq.CommonReq = req.CommonReq
	if saveResultErr := user.SaveLoginLogAsync(ctx, loginLogReq); saveResultErr != nil {
		global.Logger.Error("保存日志异常", zap.Any("err", saveResultErr), zap.Any("data", loginLogReq))
		return
	}
	baseAuthResp, baseAuthErr := user.BuildUserAuthResp(ctx, serviceResp)
	if baseAuthErr != nil {
		err = baseAuthErr
		global.Logger.Error("构建授权返回异常", zap.Any("err", baseAuthErr), zap.Any("data", serviceResp))
		return
	}
	resp.BaseUserAuthRespModel = baseAuthResp
	return
}

// AccountLogin 账号登录
func (receiver *CommonLogic) AccountLogin(ctx context.Context, req *api.AccountLoginReq) (resp api.AccountLoginResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	userService := &user.UserService{}
	serviceResp, serviceErr := userService.AccountLogin(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("账号登录异常", zap.Any("err", serviceErr))
		return
	}

	resp.PopUp = (&pop_up.PopUpService{}).GetPopUpConfig(ctx, api.BuildPopUpReq{CommonReq: req.CommonReq, UserId: serviceResp.Id, Source: api.BuildPopUpSourceLogin})

	//写登录日志
	loginLogReq := &api.LoginLogReq{}
	loginLogReq.UserId = serviceResp.Id
	loginLogReq.CommonReq = req.CommonReq
	if saveResultErr := user.SaveLoginLogAsync(ctx, loginLogReq); saveResultErr != nil {
		global.Logger.ErrorCtx(ctx, "保存日志异常", zap.Any("err", saveResultErr), zap.Any("data", loginLogReq))
		return
	}

	authData, authErr := user.BuildUserAuthResp(ctx, serviceResp)
	if authErr != nil {
		err = authErr
		global.Logger.Error("构建授权异常", zap.Any("err", authErr), zap.Any("data", serviceResp))
		return
	}
	resp.BaseUserAuthRespModel = authData
	return
}

// SmsSend 短信发送
func (receiver *CommonLogic) SmsSend(ctx context.Context, req *api.SmsSendReq) (resp api.SmsSendResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	userService := &user.UserService{}
	serviceResp, serviceErr := userService.SmsSend(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("短信发送异常", zap.Any("err", serviceErr))
		return
	}

	resp = serviceResp
	return
}

// SmsLogin 短信验证码登录
func (receiver *CommonLogic) SmsLogin(ctx context.Context, req *api.SmsLoginReq) (resp api.SmsLoginResp, err error) {
	userService := &user.UserService{}
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	serviceResp, serviceErr := userService.SmsLogin(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("短信登录异常", zap.Error(serviceErr))
		return
	}

	resp.PopUp = (&pop_up.PopUpService{}).GetPopUpConfig(ctx, api.BuildPopUpReq{CommonReq: req.CommonReq, UserId: serviceResp.Id, Source: api.BuildPopUpSourceLogin})

	//写登录日志
	loginLogReq := &api.LoginLogReq{}
	loginLogReq.Auth.UserId = serviceResp.Id
	loginLogReq.CommonReq = req.CommonReq
	if saveResultErr := user.SaveLoginLogAsync(ctx, loginLogReq); saveResultErr != nil {
		global.Logger.ErrorCtx(ctx, "保存日志异常", zap.Any("err", saveResultErr), zap.Any("data", loginLogReq))
		return
	}
	authData, authErr := user.BuildUserAuthResp(ctx, serviceResp)
	if authErr != nil {
		err = authErr
		global.Logger.ErrorCtx(ctx, "构建授权异常", zap.Any("err", authErr), zap.Any("data", serviceResp))
		return
	}
	resp.BaseUserAuthRespModel = authData
	return
}

// TokenLogin TOKEN登录
func (receiver *CommonLogic) TokenLogin(ctx context.Context, req *api.TokenLoginReq) (resp api.TokenLoginResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	userService := &user.UserService{}
	serviceResp, serviceErr := userService.TokenLogin(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("登录异常", zap.Any("err", serviceErr))
		return
	}

	resp.PopUp = (&pop_up.PopUpService{}).GetPopUpConfig(ctx, api.BuildPopUpReq{CommonReq: req.CommonReq, UserId: serviceResp.Id, Source: api.BuildPopUpSourceLogin})

	//写登录日志
	loginLogReq := &api.LoginLogReq{}
	loginLogReq.UserId = serviceResp.Id
	loginLogReq.CommonReq = req.CommonReq
	if saveResultErr := user.SaveLoginLogAsync(ctx, loginLogReq); saveResultErr != nil {
		global.Logger.ErrorCtx(ctx, "保存日志异常", zap.Any("err", saveResultErr), zap.Any("data", loginLogReq))
		return
	}
	authData, authErr := user.BuildUserAuthResp(ctx, serviceResp)
	if authErr != nil {
		err = authErr
		global.Logger.Error("构建授权异常", zap.Error(authErr), zap.Any("data", serviceResp))
		return
	}
	resp.BaseUserAuthRespModel = authData
	return
}

// PasswordRetrievePhone 通过手机找回密码
func (receiver *CommonLogic) PasswordRetrievePhone(ctx context.Context, req *api.PasswordRetrievePhoneReq) (
	resp api.PasswordRetrievePhoneResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	userService := &user.UserService{}
	serviceResp, serviceErr := userService.PasswordRetrievePhone(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("登录异常", zap.Any("err", serviceErr))
		return
	}
	resp = serviceResp
	return
}

// AccountRetrievePhone 通过手机找回账号
func (receiver *CommonLogic) AccountRetrievePhone(ctx context.Context, req *api.AccountRetrievePhoneReq) (resp interface{}, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	userService := &user.UserService{}
	serviceResp, serviceErr := userService.AccountRetrievePhone(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("登录异常", zap.Any("err", serviceErr))
		return
	}
	resp = serviceResp
	return
}

// Pay 支付
func (receiver *CommonLogic) Pay(ctx context.Context, req *api.PayReq) (resp api.PayResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if common.GetPayTypeName(req.PayType) == "" {
		err = errors2.New("支付方式未知" + req.PayType)
		return
	}

	popUpConfig := (&pop_up.PopUpService{}).GetPopUpConfig(ctx, api.BuildPopUpReq{CommonReq: req.CommonReq, UserId: req.UserId, Source: api.BuildPopUpSourcePay})
	resp.PopUp = popUpConfig

	userService := &user.UserService{}
	serviceResp, serviceErr := userService.Pay(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("登录异常", zap.Any("err", serviceErr))
		return
	}
	resp = serviceResp
	resp.PopUp = popUpConfig
	return
}

// PaymentStatusInquiry 支付状态查询
func (receiver *CommonLogic) PaymentStatusInquiry(ctx context.Context, req *api.PaymentStatusInquiryReq) (
	resp api.PaymentStatusInquiryResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	userService := &user.UserService{}
	serviceResp, serviceErr := userService.PaymentStatusInquiry(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("登录异常", zap.Any("err", serviceErr))
		return
	}
	resp = serviceResp
	return
}

// PasswordModify 密码修改
func (receiver *CommonLogic) PasswordModify(ctx context.Context, req *api.PasswordModifyReq) (
	resp api.PasswordModifyResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	userService := &user.UserService{}
	serviceResp, serviceErr := userService.PasswordModify(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("登录异常", zap.Any("err", serviceErr))
		return
	}
	resp = serviceResp
	return
}

// PhoneBind 手机绑定
func (receiver *CommonLogic) PhoneBind(ctx context.Context, req *api.PhoneBindReq) (
	resp api.PhoneBindResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	userService := &user.UserService{}
	serviceResp, serviceErr := userService.PhoneBind(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("登录异常", zap.Any("err", serviceErr))
		return
	}
	resp = serviceResp
	return
}

// IdCardBind 手机绑定
func (receiver *CommonLogic) IdCardBind(ctx context.Context, req *api.IdCardBindReq) (
	resp api.IdCardBindResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	userService := &user.UserService{}
	serviceResp, serviceErr := userService.IdCardBind(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("登录异常", zap.Any("err", serviceErr))
		return
	}
	resp = serviceResp
	return
}

// AccountDelete 用户注销
func (receiver *CommonLogic) AccountDelete(ctx context.Context, req *api.AccountDeleteReq) (
	resp api.AccountDeleteResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	userService := &user.UserService{}
	serviceResp, serviceErr := userService.AccountStatusModify(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("登录异常", zap.Any("err", serviceErr))
		return
	}
	resp = serviceResp
	return
}

func (receiver *CommonLogic) AccountLogout(ctx context.Context, req *api.AccountLogoutReq) (
	resp api.AccountLogoutResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	userService := &user.UserService{}
	serviceResp, serviceErr := userService.AccountLogout(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("异常", zap.Any("err", serviceErr))
		return
	}
	resp = serviceResp
	return
}

func (receiver *CommonLogic) Heartbeat(ctx context.Context, req *api.HeartbeatReq) (
	resp api.HeartbeatResp, err error) {
	if validateErr := req.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}

	popUpConfig := (&pop_up.PopUpService{}).GetPopUpConfig(ctx, api.BuildPopUpReq{CommonReq: req.CommonReq, UserId: req.UserId, Source: api.BuildPopUpSourceHeartbeat})
	resp.PopUp = popUpConfig

	userService := &user.UserService{}
	_, serviceErr := userService.Heartbeat(ctx, req)
	if serviceErr != nil {
		err = serviceErr
		global.Logger.Error("异常", zap.Any("err", serviceErr))
		return
	}
	return
}
