package user

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/internal/service/pay"
	"cngamesdk.com/api/internal/service/pay/channel"
	sms2 "cngamesdk.com/api/internal/service/sms"
	"cngamesdk.com/api/internal/service/token"
	"cngamesdk.com/api/model/api"
	error3 "cngamesdk.com/api/model/api/error"
	"cngamesdk.com/api/model/sql/log"
	"cngamesdk.com/api/model/sql/user"
	"context"
	"errors"
	error2 "github.com/cngamesdk/go-core/model/error"
	user2 "github.com/cngamesdk/go-core/model/sql/user"
	"github.com/cngamesdk/go-core/util/identity"
	"github.com/cngamesdk/go-core/validate"
	"github.com/duke-git/lancet/v2/random"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/panjf2000/ants/v2"
	errors2 "github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

var (
	saveLoginLogPool *ants.Pool
)

func init() {
	mySaveLoginLogPool, mySaveLoginLogPoolErr := ants.NewPool(100)
	if mySaveLoginLogPoolErr != nil {
		global.Logger.Error("创建协程池异常", zap.Any("err", mySaveLoginLogPoolErr))
		return
	}
	saveLoginLogPool = mySaveLoginLogPool
}

type UserService struct {
}

// AccountAllocate 账号分配
func (receiver *UserService) AccountAllocate(ctx context.Context, req *api.AccountAllocateReq) (resp api.AccountAllocateResp, err error) {
	resp.UserName = getAccountPrex() + random.RandNumeral(6)
	resp.Password = random.RandNumeral(6)
	return
}

// userNameAvailable 验证用户名是否可用
func (receiver *UserService) userNameAvailable(ctx context.Context, platformId int64, userName string) (resp *user.OdsUserInfoLogModel, err error) {
	if validateErr := validate.EmptyString(userName); validateErr != nil {
		err = validateErr
		return
	}
	userModel := user.NewOdsUserInfoLogModel()
	takeErr := userModel.Take(ctx, "*", "platform_id = ? and user_name = ?", platformId, userName)
	if takeErr == nil {
		err = errors2.Wrap(error2.ErrorRecordIsExists, "用户名")
		resp = userModel
		return
	}
	if !errors.Is(takeErr, error2.ErrorRecordIsNotFind) {
		err = takeErr
		global.Logger.Error("获取用户信息出错", zap.String("err", takeErr.Error()))
		return
	}
	return
}

// AccountReg 账号注册
func (receiver *UserService) AccountReg(ctx context.Context, req *api.AccountRegReq) (resp *user.OdsUserInfoLogModel, err error) {
	//验证用户名是否已经存在
	if _, checkErr := receiver.userNameAvailable(ctx, req.PlatformId, req.UserName); checkErr != nil {
		err = checkErr
		global.Logger.Error("验证用户名存在异常", zap.Error(checkErr))
		return
	}
	userModel := user.NewOdsUserInfoLogModel()
	userModel.PlatformId = req.PlatformId
	userModel.UserName = req.UserName
	userModel.Password = req.Password
	userModel.TrueName = req.TrueName
	userModel.IdCard = req.IdCard
	userModel.Phone = req.Phone
	userModel.Status = user2.UserStatusNormal
	if createErr := userModel.Create(ctx); createErr != nil {
		err = createErr
		global.Logger.Error("保存用户异常", zap.String("err", createErr.Error()), zap.Any("data", userModel))
		return
	}

	if saveRegLogErr := saveRegLog(ctx, userModel, req.CommonReq); saveRegLogErr != nil {
		err = saveRegLogErr
		global.Logger.Error("保存用户日志异常", zap.Any("err", saveRegLogErr))
		return
	}

	resp = userModel
	return
}

// AccountLogin 账号登录
func (receiver *UserService) AccountLogin(ctx context.Context, req *api.AccountLoginReq) (resp *user.OdsUserInfoLogModel, err error) {
	userModel := user.NewOdsUserInfoLogModel()
	if takeErr := userModel.Take(ctx, "*", "platform_id = ? and user_name = ?", req.PlatformId, req.UserName); takeErr != nil {
		err = error3.ErrorUserNameOrPasswordValid
		global.Logger.ErrorCtx(ctx, "获取用户异常", zap.Error(takeErr), zap.Any("data", req))
		return
	}
	if !userModel.ValidatePassword(req.Password) {
		err = error3.ErrorUserNameOrPasswordValid
		return
	}
	if !userModel.Valid() {
		err = error3.ErrorUserStatusValid
		global.Logger.ErrorCtx(ctx, "用户状态异常", zap.Any("err", err), zap.Any("data", userModel.Status))
		return
	}
	resp = userModel
	return
}

// SmsSend 短信发送
func (receiver *UserService) SmsSend(ctx context.Context, req *api.SmsSendReq) (resp api.SmsSendResp, err error) {
	smsService := sms2.NewSmsService(req.Phone)
	sendErr := smsService.SendCheckCode(ctx, req)
	if sendErr != nil {
		err = sendErr
		global.Logger.Error("发送验证码异常", zap.Any("err", sendErr), zap.Any("data", req))
		return
	}
	return
}

// SmsLogin 短信登录
func (receiver *UserService) SmsLogin(ctx context.Context, req *api.SmsLoginReq) (resp *user.OdsUserInfoLogModel, err error) {
	smsService := sms2.NewSmsService(req.Phone)
	if validateErr := smsService.VerifyCheckCode(ctx, req.Code); validateErr != nil {
		err = validateErr
		return
	}
	userModel, availableErr := receiver.userNameAvailable(ctx, req.PlatformId, req.Phone)
	if availableErr != nil {
		if !errors.Is(availableErr, error2.ErrorRecordIsExists) {
			err = error2.ErrorInternalSystem
			global.Logger.Error("验证用户名存在异常", zap.Any("err", availableErr))
			return
		}
	}
	//新用户
	if userModel == nil {
		userModel = user.NewOdsUserInfoLogModel()
		userModel.PlatformId = req.PlatformId
		userModel.UserName = req.Phone
		userModel.Password = random.RandNumeral(6)
		userModel.Phone = req.Phone
		userModel.Status = user2.UserStatusNormal
		if createErr := userModel.Create(ctx); createErr != nil {
			err = createErr
			global.Logger.Error("保存用户异常", zap.String("err", createErr.Error()), zap.Any("data", userModel))
			return
		}
		if saveRegLogErr := saveRegLog(ctx, userModel, req.CommonReq); saveRegLogErr != nil {
			err = saveRegLogErr
			global.Logger.ErrorCtx(ctx, "保存注册用户日志异常", zap.Any("err", saveRegLogErr), zap.Any("data", userModel))
			return
		}
	}
	if !userModel.Valid() {
		err = error3.ErrorUserStatusValid
		global.Logger.ErrorCtx(ctx, "用户状态异常", zap.Any("err", err), zap.Any("data", userModel.Status))
		return
	}
	resp = userModel
	return
}

// TokenLogin TOKEN登录
func (receiver *UserService) TokenLogin(ctx context.Context, req *api.TokenLoginReq) (resp *user.OdsUserInfoLogModel, err error) {
	tokenService := &token.TokenService{}
	_, tokenErr := tokenService.Verify(ctx, req.Token)
	if tokenErr != nil {
		err = tokenErr
		global.Logger.Error("TOKEN验证异常", zap.Error(tokenErr))
		return
	}
	userModel := tokenService.GetUserModel()
	resp = userModel
	return
}

// PasswordRetrievePhone 通过手机找回密码
func (receiver *UserService) PasswordRetrievePhone(ctx context.Context, req *api.PasswordRetrievePhoneReq) (
	resp api.PasswordRetrievePhoneResp, err error) {
	smsService := sms2.NewSmsService(req.Phone)
	if validateErr := smsService.VerifyCheckCode(ctx, req.Code); validateErr != nil {
		err = validateErr
		return
	}
	userModel := user.NewOdsUserInfoLogModel()
	if takeErr := userModel.Take(ctx, "id,phone_crypt", "platform_id = ? and user_name = ?", req.PlatformId, req.UserName); takeErr != nil {
		err = takeErr
		global.Logger.Error("获取用户异常", zap.Any("err", takeErr), zap.Any("data", req.UserName))
		return
	}
	if userModel.Phone != req.Phone {
		err = errors2.New("手机号与账号绑定不一致")
		global.Logger.Warn("手机号与账号绑定不一致", zap.Any("err", err))
		return
	}
	password := req.Password
	updateUserModel := user.NewOdsUserInfoLogModel()
	updateUserModel.Id = userModel.Id
	updateUserModel.PlatformId = req.PlatformId
	updateUserModel.Password = password
	if updateErr := updateUserModel.Updates(ctx, "id = ?", userModel.Id); updateErr != nil {
		err = error2.ErrorInternalSystem
		global.Logger.Error("更新密码异常", zap.Any("err", updateErr))
		return
	}
	return
}

// AccountRetrievePhone 通过手机找回账号
func (receiver *UserService) AccountRetrievePhone(ctx context.Context, req *api.AccountRetrievePhoneReq) (resp interface{}, err error) {
	smsService := sms2.NewSmsService(req.Phone)
	if validateErr := smsService.VerifyCheckCode(ctx, req.Code); validateErr != nil {
		err = validateErr
		return
	}
	var list []user.OdsUserInfoLogModel
	userModel := user.NewOdsUserInfoLogModel()
	phoneEncrypt := userModel.GetEcbEncrypt(req.Phone)
	if findErr := userModel.Db().
		WithContext(ctx).
		Table(userModel.TableName()).
		Select("user_name").
		Where("platform_id = ? and phone_crypt = ?", req.PlatformId, phoneEncrypt).
		Limit(5).
		Find(&list).Error; findErr != nil {
		err = findErr
		global.Logger.Error("获取列表异常", zap.Any("err", findErr), zap.Any("data", phoneEncrypt))
		return
	}
	var respData []api.AccountRetrievePhoneResp
	for _, item := range list {
		respData = append(respData, api.AccountRetrievePhoneResp{
			UserName: item.UserName,
		})
	}
	resp = respData
	return
}

// Pay SDK支付
func (receiver *UserService) Pay(ctx context.Context, req *api.PayReq) (resp api.PayResp, err error) {
	payChannel, payChannelErr := pay.GetSdkPayChannelFactory(req)
	preOrderReq := channel.PreOrderReq{
		Money: req.Money,
	}
	if payChannelErr != nil {
		err = errors2.New("获取支付渠道异常")
		global.Logger.Error("获取支付渠道异常", zap.Any("err", payChannelErr))
		return
	}
	if payChannel == nil {
		err = errors2.New("未获取到支付渠道")
		global.Logger.Error("未获取到支付渠道", zap.Any("data", req))
		return
	}

	preOrderResp, preOrderErr := payChannel.PreOrder(ctx, preOrderReq)
	if preOrderErr != nil {
		err = preOrderErr
		global.Logger.Error("预下单异常", zap.Any("err", preOrderErr))
		return
	}

	payLogReq := api.PayLogReq{}
	payLogReq.PayReq = *req
	payLogReq.PayChannelId = payChannel.GetPayChannelId()
	payLogReq.OrderId = preOrderResp.OrderId
	payLogReq.CommonReq = req.CommonReq
	if _, saveErr := pay.SavePayLog(ctx, payLogReq); saveErr != nil {
		err = saveErr
		global.Logger.Error("保存日志异常", zap.Any("err", saveErr), zap.Any("data", payLogReq))
		return
	}

	resp.OrderId = preOrderResp.OrderId
	resp.Url = preOrderResp.Url
	return
}

// PaymentStatusInquiry 支付状态查询
func (receiver *UserService) PaymentStatusInquiry(ctx context.Context, req *api.PaymentStatusInquiryReq) (
	resp api.PaymentStatusInquiryResp, err error) {
	payModel := log.NewOdsPayLogModel()
	if takeErr := payModel.Take(ctx, "pay_status,user_id", "platform_id = ? and order_id = ?", req.PlatformId, req.OrderId); takeErr != nil {
		err = takeErr
		global.Logger.Error("获取异常", zap.Any("err", takeErr))
		return
	}
	if payModel.UserId != req.UserId {
		err = errors.New("非法请求")
		global.Logger.Error("下单用户与TOKEN用户不一致", zap.Any("err", err))
		return
	}
	resp.Status = payModel.PayStatus
	return
}

// PasswordModify 密码修改
func (receiver *UserService) PasswordModify(ctx context.Context, req *api.PasswordModifyReq) (
	resp api.PasswordModifyResp, err error) {
	userModel := user.NewOdsUserInfoLogModel()
	userModel.Id = req.UserId
	userModel.PlatformId = req.PlatformId
	userModel.Password = req.Password
	if saveErr := userModel.Updates(ctx, "id = ?", req.UserId); saveErr != nil {
		err = saveErr
		global.Logger.Error("修改密码异常", zap.Any("err", saveErr))
		return
	}
	return
}

// PhoneBind 手机绑定
func (receiver *UserService) PhoneBind(ctx context.Context, req *api.PhoneBindReq) (
	resp api.PhoneBindResp, err error) {
	smsService := sms2.NewSmsService(req.Phone)
	if validateErr := smsService.VerifyCheckCode(ctx, req.Code); validateErr != nil {
		err = validateErr
		return
	}
	userModel := user.NewOdsUserInfoLogModel()
	userModel.Id = req.UserId
	userModel.PlatformId = req.PlatformId
	userModel.Phone = req.Phone
	if saveErr := userModel.Updates(ctx, "id = ?", req.UserId); saveErr != nil {
		err = saveErr
		global.Logger.Error("修改异常", zap.Any("err", saveErr))
		return
	}
	return
}

// IdCardBind 身份证绑定
func (receiver *UserService) IdCardBind(ctx context.Context, req *api.IdCardBindReq) (
	resp api.IdCardBindResp, err error) {
	userModel := user.NewOdsUserInfoLogModel()
	userModel.Id = req.UserId
	userModel.PlatformId = req.PlatformId
	userModel.IdCard = req.IdCard
	userModel.TrueName = req.TrueName
	if saveErr := userModel.Updates(ctx, "id = ?", req.UserId); saveErr != nil {
		err = saveErr
		global.Logger.Error("修改异常", zap.Any("err", saveErr))
		return
	}
	return
}

// AccountStatusModify 用户状态更改
func (receiver *UserService) AccountStatusModify(ctx context.Context, req *api.AccountDeleteReq) (
	resp api.AccountDeleteResp, err error) {
	userModel := user.NewOdsUserInfoLogModel()
	userModel.Id = req.UserId
	userModel.PlatformId = req.PlatformId
	userModel.Status = req.Status
	if saveErr := userModel.Updates(ctx, "id = ?", req.UserId); saveErr != nil {
		err = saveErr
		global.Logger.Error("修改异常", zap.Any("err", saveErr))
		return
	}
	return
}

// AccountLogout 用户账号登出
func (receiver *UserService) AccountLogout(ctx context.Context, req *api.AccountLogoutReq) (resp api.AccountLogoutResp, err error) {
	tokenService := &token.TokenService{UserId: req.UserId}
	err = tokenService.DelCacheToken(ctx)
	if err != nil {
		global.Logger.Error("删除缓存异常", zap.Error(err))
		return
	}
	return
}

// Heartbeat 心跳
func (receiver *UserService) Heartbeat(ctx context.Context, req *api.HeartbeatReq) (resp api.HeartbeatResp, err error) {
	return
}

// saveRegLog 保存注册日志
func saveRegLog(ctx context.Context, userModel *user.OdsUserInfoLogModel, commonReq api.CommonReq) (err error) {
	logModel := log.NewOdsRegLogModel()
	logModel.PlatformId = userModel.PlatformId
	logModel.UserId = userModel.Id
	logModel.Ipv6 = commonReq.Ipv6
	logModel.Ipv4 = commonReq.Ipv4
	logModel.AndriodId = commonReq.AndriodId
	logModel.ClientIp = commonReq.ClientIp
	logModel.Idfv = commonReq.Idfv
	logModel.Imei = commonReq.Imei
	logModel.Network = commonReq.Network
	logModel.Oaid = commonReq.Oaid
	logModel.SystemVersion = commonReq.SystemVersion
	logModel.SiteId = commonReq.SiteId
	logModel.AgentId = commonReq.AgentId
	logModel.RegTime = time.Now()
	logModel.MediaSiteId = commonReq.MediaSiteId
	logModel.Model = commonReq.Model
	logModel.GameId = commonReq.GameId
	logModel.Brand = commonReq.Brand
	logModel.ChannelId = commonReq.ChannelId
	logModel.UserAgent = commonReq.UserAgent

	//交换存储
	if commonReq.MediaSiteId > 0 {
		logModel.MediaSiteId = commonReq.SiteId
		logModel.SiteId = commonReq.MediaSiteId
	}

	if createErr := logModel.Create(ctx); createErr != nil {
		err = createErr
		global.Logger.Error("保存用户日志异常", zap.Any("err", createErr.Error()), zap.Any("data", logModel))
		return
	}
	return
}

// SaveLoginLogAsync 异步保存登录日志
func SaveLoginLogAsync(ctx context.Context, req *api.LoginLogReq) (err error) {
	err = saveLoginLogPool.Submit(func() {
		saveErr := SaveLoginLog(ctx, req)
		if saveErr != nil {
			global.Logger.ErrorCtx(ctx, "保存日志异常", zap.Any("err", saveErr), zap.Any("data", req))
			return
		}
	})
	return
}

// SaveLoginLog 保存登录日志
func SaveLoginLog(ctx context.Context, req *api.LoginLogReq) (err error) {
	logModel := log.NewOdsLoginLogModel()
	logModel.PlatformId = req.PlatformId
	logModel.UserId = req.UserId
	logModel.Ipv6 = req.Ipv6
	logModel.Ipv4 = req.Ipv4
	logModel.AndriodId = req.AndriodId
	logModel.ClientIp = req.ClientIp
	logModel.Idfv = req.Idfv
	logModel.Imei = req.Imei
	logModel.Network = req.Network
	logModel.Oaid = req.Oaid
	logModel.SystemVersion = req.SystemVersion
	logModel.SiteId = req.SiteId
	logModel.AgentId = req.AgentId
	logModel.LoginTime = time.Now()
	logModel.MediaSiteId = req.MediaSiteId
	logModel.Model = req.Model
	logModel.GameId = req.GameId
	logModel.Brand = req.Brand
	logModel.ChannelId = req.ChannelId
	logModel.UserAgent = req.UserAgent

	if req.MediaSiteId > 0 {
		logModel.MediaSiteId = req.SiteId
		logModel.SiteId = req.MediaSiteId
	}

	if createErr := logModel.Create(ctx); createErr != nil {
		err = createErr
		global.Logger.ErrorCtx(ctx, "保存用户日志异常", zap.Any("err", createErr.Error()), zap.Any("data", logModel))
		return
	}
	return
}

// getAccountPrex 获取随机账号前缀
func getAccountPrex() string {
	accountPrexsLen := len(global.Config.Common.AccountPrefixs)
	if accountPrexsLen <= 0 {
		return ""
	}
	randomIndex := random.RandInt(0, accountPrexsLen)
	return global.Config.Common.AccountPrefixs[randomIndex]
}

// BuildUserAuthResp 构建用户登录授权后返回
func BuildUserAuthResp(ctx context.Context, req *user.OdsUserInfoLogModel) (resp api.BaseUserAuthRespModel, err error) {
	resp.UserId = req.Id
	resp.Phone = strutil.HideString(req.Phone, 3, 7, "*")

	if req.IdCard != "" {

		resp.RealName = 1

		idCardInstance := identity.New(req.IdCard)
		idCardParseErr := idCardInstance.Parse()
		if idCardParseErr != nil {
			err = idCardParseErr
			global.Logger.ErrorCtx(ctx, "身份证解析异常", zap.Any("err", idCardParseErr), zap.Any("data", req.IdCard))
			return
		}
		resp.Age = identity.Age(idCardInstance.GetBirthdayTime())
	}
	tokenService := token.NewTokenService(req.Id, req.Version)
	tokenResp, tokenErr := tokenService.Generate(ctx, nil)
	if tokenErr != nil {
		err = tokenErr
		global.Logger.ErrorCtx(ctx, "生成token异常", zap.Any("err", tokenErr), zap.Any("data", req.Id))
		return
	}
	resp.Token = tokenResp
	return
}

// BuildPopUp 构建弹窗
func BuildPopUp(ctx context.Context, req api.BuildPopUpReq) (resp api.PopUpConfig, err error) {
	if req.UserId <= 0 {
		resp = api.PopUpConfig{
			Show: 1,
			Url:  "https://www.baidu.com/",
			Btns: []api.PopUpConfigBtn{
				{Type: api.PopUpConfigBtnConfirm, Text: "好的"},
			},
		}
	} else {
		resp = api.PopUpConfig{
			Show: 1,
			Url:  "https://www.baidu.com/",
			Btns: []api.PopUpConfigBtn{
				{Type: api.PopUpConfigBtnCancel, Text: "取消"},
				{Type: api.PopUpConfigBtnConfirm, Text: "确认"},
			},
		}
	}
	return
}
