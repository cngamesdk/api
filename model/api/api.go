package api

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/model/cache/game"
	"cngamesdk.com/api/model/sql/log"
	"context"
	"fmt"
	error2 "github.com/cngamesdk/go-core/model/error"
	"github.com/cngamesdk/go-core/model/sql/common"
	user2 "github.com/cngamesdk/go-core/model/sql/user"
	"github.com/cngamesdk/go-core/validate"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/duke-git/lancet/v2/validator"
	"github.com/gin-gonic/gin"
	errors2 "github.com/pkg/errors"
	"github.com/spf13/cast"
	"strings"
)

type CommonReq struct {
	Platform       string `json:"platform" form:"platform" binding:"required"` // 对平台ID AES加密后值
	PlatformId     int64  `json:"platform_id" form:"platform_id"`
	GameId         int64  `json:"game_id" form:"game_id"`
	AgentId        int64  `json:"agent_id" form:"agent_id"`
	SiteId         int64  `json:"site_id" form:"site_id"`
	Idfv           string `json:"idfv" form:"idfv"`
	Imei           string `json:"imei" form:"imei"`
	AndriodId      string `json:"andriod_id" form:"andriod_id"`
	Oaid           string `json:"oaid" form:"oaid"`
	Brand          string `json:"brand" form:"brand"`
	Model          string `json:"model" form:"model"`
	SystemVersion  string `json:"system_version" form:"system_version"`
	Network        string `json:"network" form:"network"`
	AppVersionCode int    `json:"app_version_code" form:"app_version_code"`
	SdkVersionCode int    `json:"sdk_version_code" form:"sdk_version_code"`
	MediaSiteId    int64  `json:"media_site_id" form:"media_site_id"`
	Ipv4           string `json:"ipv4" form:"ipv4"`
	Ipv6           string `json:"ipv6" form:"ipv6"`
	ClientIp       string `json:"client_ip" form:"client_ip"`
	ChannelId      int64  `json:"channel_id" form:"channel_id"`
	UserAgent      string `json:"user_agent" form:"user_agent"`
	Media          string `json:"media" form:"media"`
	Timestamp      int64  `json:"timestamp" form:"timestamp"`
	Sign           string `json:"sign" form:"sign"`
}

func (receiver *CommonReq) Format(ctx context.Context) {
	receiver.PlatformId = 0
	if receiver.Platform != "" {
		receiver.PlatformId = cast.ToInt64(string(cryptor.AesEcbDecrypt([]byte(cryptor.Base64StdDecode(receiver.Platform)), []byte(global.Config.Common.AesCryptKey))))
	}
	if receiver.ClientIp == "" {
		ctxGin, ok := ctx.(*gin.Context)
		if ok {
			receiver.ClientIp = ctxGin.ClientIP()
			if validator.IsIpV6(receiver.ClientIp) {
				if receiver.Ipv6 == "" {
					receiver.Ipv6 = receiver.ClientIp
				}
			} else if receiver.Ipv4 == "" {
				receiver.Ipv4 = receiver.ClientIp
			}
			if receiver.UserAgent == "" {
				receiver.UserAgent = ctxGin.GetHeader("User-Agent")
			}
		}
	}
	receiver.Sign = strings.TrimSpace(receiver.Sign)
}

func (receiver *CommonReq) Validate(ctx context.Context) (err error) {
	ctxClient := ctx.Value(global.CtxKeyClient)
	client, clientOk := ctxClient.(string)
	if clientOk {
		if client == common.GameTypeMobileGame {
			if validateErr := receiver.validateMobileGame(ctx); validateErr != nil {
				err = validateErr
				return
			}
		}
	}
	if validateErr := validate.EmptyString(receiver.Platform); validateErr != nil {
		err = errors2.Wrap(error2.ErrorParamEmpty, "platform")
		return
	}
	if receiver.PlatformId <= 0 {
		err = errors2.Wrap(error2.ErrorParamEmpty, "platform解密失败")
		return
	}
	return
}

// sdkCommonValidate SDK通用验证
func (receiver *CommonReq) sdkCommonValidate(ctx context.Context) (err error) {
	if receiver.GameId <= 0 {
		err = errors2.Wrap(error2.ErrorParamEmpty, "game_id")
		return
	}
	appKey := common.GetGameAppKey(receiver.GameId, global.Config.Common.GameHashKey)
	if cryptor.Md5String(fmt.Sprintf("%d%s", receiver.Timestamp, appKey)) != receiver.Sign {
		err = errors2.Wrap(error2.ErrorSignVerifyFail, "sign")
		return
	}
	if validateErr := receiver.validatePlatformAndGameMap(ctx); validateErr != nil {
		err = validateErr
		return
	}
	return
}

// validateMobileGame 验证手游
func (receiver *CommonReq) validateMobileGame(ctx context.Context) (err error) {
	if validateErr := receiver.sdkCommonValidate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	return
}

// validatePlatformAndGameMap 验证平台和游戏对应关系
func (receiver *CommonReq) validatePlatformAndGameMap(ctx context.Context) (err error) {
	model := game.NewDimGameModel()
	if takeErr := model.Take(ctx, "*", "id = ?", receiver.GameId); takeErr != nil {
		err = takeErr
		return
	}
	if model.PlatformId != receiver.PlatformId {
		err = errors2.New("platform非法请求")
		return
	}
	return
}

type LoginLogReq struct {
	Auth
	CommonReq
}

type LoginLogResp struct {
}

type BaseUserAuthRespModel struct {
	Token    string `json:"token"`
	UserId   int64  `json:"user_id"`
	RealName int    `json:"real_name"`
	Phone    string `json:"phone"`
	Age      int    `json:"age"`
}

type InitReq struct {
	CommonReq
}

func (receiver *InitReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
}

func (receiver *InitReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	return
}

type InitResp struct {
	QqNumber       string          `json:"qq_number"`
	QqLink         string          `json:"qq_link"`
	DomainList     []string        `json:"domain_list"`
	PayConfig      []string        `json:"pay_config"`
	PopUp          PopUpConfig     `json:"pop_up"`
	Heartbeat      HeartbeatSwitch `json:"heartbeat"`
	SwitchRealName int             `json:"switch_real_name"`
}

type AccountAllocateReq struct {
	CommonReq
}

func (receiver *AccountAllocateReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
}

func (receiver *AccountAllocateReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	return
}

type AccountAllocateResp struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type AccountRegReq struct {
	CommonReq
	UserName string `json:"user_name" form:"user_name" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Phone    string `json:"phone" form:"phone"`
	TrueName string `json:"true_name" form:"true_name"`
	IdCard   string `json:"id_card" form:"id_card"`
}

func (receiver *AccountRegReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.UserName = strings.ToLower(strings.TrimSpace(receiver.UserName))
	receiver.Password = strings.TrimSpace(receiver.Password)
	receiver.Phone = strings.TrimSpace(receiver.Phone)
	receiver.TrueName = strings.TrimSpace(receiver.TrueName)
	receiver.IdCard = strings.TrimSpace(receiver.IdCard)
}

func (receiver *AccountRegReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.UserName(receiver.UserName); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.Password(receiver.Password); validateErr != nil {
		err = validateErr
		return
	}
	if receiver.Phone != "" {
		if validateErr := validate.ChineseMobile(receiver.Phone); validateErr != nil {
			err = validateErr
			return
		}
	}
	if receiver.IdCard != "" {
		if validateErr := validate.ChineseIdCard(receiver.IdCard); validateErr != nil {
			err = validateErr
			return
		}
		if validateErr := validate.EmptyString(receiver.TrueName); validateErr != nil {
			err = validateErr
			return
		}
	}

	return
}

type AccountRegResp struct {
	BaseUserAuthRespModel
}

type AccountLoginReq struct {
	CommonReq
	UserName string `json:"user_name" form:"user_name" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

func (receiver *AccountLoginReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.UserName = strings.ToLower(strings.TrimSpace(receiver.UserName))
	receiver.Password = strings.TrimSpace(receiver.Password)
}

func (receiver *AccountLoginReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.UserName(receiver.UserName); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.Password(receiver.Password); validateErr != nil {
		err = validateErr
		return
	}
	return
}

type PopUpResp struct {
	PopUp PopUpConfig `json:"pop_up"`
}

type AccountLoginResp struct {
	BaseUserAuthRespModel
	PopUpResp
}

type SmsSendReq struct {
	CommonReq
	Phone string `json:"phone" form:"phone" binding:"required"`
}

func (receiver *SmsSendReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.Phone = strings.ToLower(strings.TrimSpace(receiver.Phone))
}

func (receiver *SmsSendReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.ChineseMobile(receiver.Phone); validateErr != nil {
		err = validateErr
		return
	}
	return
}

type SmsSendResp struct {
}

type SmsLoginReq struct {
	CommonReq
	Phone string
	Code  string
}

func (receiver *SmsLoginReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.Phone = strings.ToLower(strings.TrimSpace(receiver.Phone))
	receiver.Code = strings.ToLower(strings.TrimSpace(receiver.Code))
}

func (receiver *SmsLoginReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.ChineseMobile(receiver.Phone); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.EmptyString(receiver.Code); validateErr != nil {
		err = errors2.Wrap(validateErr, "验证码")
		return
	}
	return
}

type SmsLoginResp struct {
	PopUpResp
	BaseUserAuthRespModel
}

type TokenLoginReq struct {
	CommonReq
	Token string `json:"token" form:"token" binding:"required"`
}

func (receiver *TokenLoginReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.Token = strings.TrimSpace(receiver.Token)
}

func (receiver *TokenLoginReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.EmptyString(receiver.Token); validateErr != nil {
		err = errors2.Wrap(validateErr, "TOKEN")
		return
	}
	return
}

type TokenLoginResp struct {
	PopUpResp
	BaseUserAuthRespModel
}

type PasswordRetrievePhoneReq struct {
	CommonReq
	UserName        string `json:"user_name" form:"user_name" binding:"required"`
	Phone           string `json:"phone" form:"phone" binding:"required"`
	Password        string `json:"password" form:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"required"`
	Code            string `json:"code" form:"code" binding:"required"`
}

func (receiver *PasswordRetrievePhoneReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.UserName = strings.ToLower(strings.TrimSpace(receiver.UserName))
	receiver.Phone = strings.ToLower(strings.TrimSpace(receiver.Phone))
	receiver.Code = strings.ToLower(strings.TrimSpace(receiver.Code))
	receiver.Password = strings.TrimSpace(receiver.Password)
	receiver.ConfirmPassword = strings.TrimSpace(receiver.ConfirmPassword)
}

func (receiver *PasswordRetrievePhoneReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.EmptyString(receiver.UserName); validateErr != nil {
		err = errors2.Wrap(validateErr, "用户名")
		return
	}
	if validateErr := validate.ChineseMobile(receiver.Phone); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.EmptyString(receiver.Code); validateErr != nil {
		err = errors2.Wrap(error2.ErrorParamEmpty, "验证码")
		return
	}
	if validateErr := validate.Password(receiver.Password); validateErr != nil {
		err = validateErr
		return
	}
	if receiver.Password != receiver.ConfirmPassword {
		err = errors2.New("密码与确认密码不一致")
		return
	}
	return
}

type PasswordRetrievePhoneResp struct {
}

type AccountRetrievePhoneReq struct {
	CommonReq
	Phone string `json:"phone" form:"phone" binding:"required"`
	Code  string `json:"code" form:"code" binding:"required"`
}

func (receiver *AccountRetrievePhoneReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.Phone = strings.ToLower(strings.TrimSpace(receiver.Phone))
}

func (receiver *AccountRetrievePhoneReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.ChineseMobile(receiver.Phone); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.EmptyString(receiver.Code); validateErr != nil {
		err = errors2.Wrap(error2.ErrorParamEmpty, "验证码")
		return
	}
	return
}

type AccountRetrievePhoneResp struct {
	UserName string `json:"user_name"`
}

type Auth struct {
	UserId int64
}

func (receiver *Auth) Format(ctx context.Context) {
	receiver.UserId = 0
	u, ok := ctx.Value(global.Config.Common.CtxTokenDataKey).(map[string]interface{})
	if ok {
		receiver.UserId = cast.ToInt64(u["user_id"])
	}
}

func (receiver *Auth) Validate(ctx context.Context) (err error) {
	if receiver.UserId <= 0 {
		err = errors2.New("用户ID为空")
		return
	}
	return
}

type PayReq struct {
	CommonReq
	PayType     string `json:"pay_type" form:"pay_type" binding:"required"`
	RoleId      string `json:"role_id" form:"role_id" binding:"required"`
	RoleName    string `json:"role_name" form:"role_name" binding:"required"`
	Money       int    `json:"money" form:"money" binding:"required"`
	ProductId   string `json:"product_id" form:"product_id" binding:"required"`
	ProductName string `json:"product_name" form:"product_name" binding:"required"`
	ServerId    int64  `json:"server_id" form:"server_id" binding:"required"`
	ServerName  string `json:"server_name" form:"server_name" binding:"required"`
	Ext         string `json:"ext" form:"ext"`
	Auth
}

func (receiver *PayReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.PayType = strings.TrimSpace(receiver.PayType)
	receiver.RoleId = strings.TrimSpace(receiver.RoleId)
	receiver.ProductId = strings.TrimSpace(receiver.ProductId)
	receiver.ProductName = strings.TrimSpace(receiver.ProductName)
	receiver.ServerName = strings.TrimSpace(receiver.ServerName)
	receiver.Ext = strings.TrimSpace(receiver.Ext)
	receiver.Auth.Format(ctx)
}

func (receiver *PayReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.EmptyString(receiver.RoleId); validateErr != nil {
		err = errors2.Wrap(validateErr, "角色ID")
		return
	}
	if validateErr := validate.EmptyString(receiver.RoleName); validateErr != nil {
		err = errors2.Wrap(validateErr, "角色名称")
		return
	}
	if receiver.Money < 0 || receiver.Money > 10000000 {
		err = errors2.New("金额无效")
		return
	}
	if validateErr := validate.EmptyString(receiver.ProductId); validateErr != nil {
		err = errors2.Wrap(validateErr, "产品ID")
		return
	}
	if validateErr := validate.EmptyString(receiver.ProductName); validateErr != nil {
		err = errors2.Wrap(validateErr, "产品名称")
		return
	}
	if receiver.ServerId <= 0 {
		err = errors2.New("区服ID为空")
		return
	}
	if validateErr := validate.EmptyString(receiver.ServerName); validateErr != nil {
		err = errors2.Wrap(validateErr, "区服名称")
		return
	}
	if validateErr := receiver.Auth.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	return
}

type PayResp struct {
	OrderId string
	Url     string
	PopUp   PopUpConfig `json:"pop_up"`
}

type PaymentStatusInquiryReq struct {
	CommonReq
	Auth
	OrderId string `json:"order_id" form:"order_id" binding:"required" `
}

func (receiver *PaymentStatusInquiryReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.OrderId = strings.TrimSpace(receiver.OrderId)
	receiver.Auth.Format(ctx)
}

func (receiver *PaymentStatusInquiryReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.EmptyString(receiver.OrderId); validateErr != nil {
		err = errors2.Wrap(validateErr, "订单号")
		return
	}
	if validateErr := receiver.Auth.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	return
}

type PaymentStatusInquiryResp struct {
	Status string `json:"status"`
}

type PasswordModifyReq struct {
	CommonReq
	Auth
	Password string `json:"password" form:"password" binding:"required"`
}

func (receiver *PasswordModifyReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.Password = strings.TrimSpace(receiver.Password)
	receiver.Auth.Format(ctx)
}

func (receiver *PasswordModifyReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.Password(receiver.Password); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := receiver.Auth.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	return
}

type PasswordModifyResp struct {
}

type PhoneBindReq struct {
	CommonReq
	Auth
	Phone string `json:"phone" form:"phone" binding:"required"`
	Code  string `json:"code" form:"code" binding:"required"`
}

func (receiver *PhoneBindReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.Phone = strings.TrimSpace(receiver.Phone)
	receiver.Code = strings.TrimSpace(receiver.Code)
	receiver.Auth.Format(ctx)
}

func (receiver *PhoneBindReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.ChineseMobile(receiver.Phone); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := receiver.Auth.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr2 := validate.EmptyString(receiver.Code); validateErr2 != nil {
		err = errors2.Wrap(validateErr2, "验证码")
		return
	}
	return
}

type PhoneBindResp struct {
}

type IdCardBindReq struct {
	CommonReq
	Auth
	IdCard   string `json:"id_card" form:"id_card" binding:"required"`
	TrueName string `json:"true_name" form:"true_name" binding:"required"`
}

func (receiver *IdCardBindReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.IdCard = strings.TrimSpace(receiver.IdCard)
	receiver.TrueName = strings.TrimSpace(receiver.TrueName)
	receiver.Auth.Format(ctx)
}

func (receiver *IdCardBindReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := receiver.Auth.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.ChineseIdCard(receiver.IdCard); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.EmptyString(receiver.TrueName); validateErr != nil {
		err = validateErr
		return
	}
	return
}

type IdCardBindResp struct {
}

type AccountDeleteReq struct {
	CommonReq
	Auth
	Status string
}

func (receiver *AccountDeleteReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.Status = user2.UserStatusDelete
	receiver.Auth.Format(ctx)
}

func (receiver *AccountDeleteReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := receiver.Auth.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if _, ok := user2.UserStatusMap[receiver.Status]; !ok {
		err = errors2.New("用户状态非法。" + receiver.Status)
		return
	}
	return
}

type AccountDeleteResp struct {
}

// AccountLogoutReq 账号登出请求
type AccountLogoutReq struct {
	CommonReq
	Auth
}

func (receiver *AccountLogoutReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.Auth.Format(ctx)
}

func (receiver *AccountLogoutReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := receiver.Auth.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	return
}

// AccountLogoutResp 账号登出返回
type AccountLogoutResp struct {
}

// HeartbeatReq 心跳请求
type HeartbeatReq struct {
	CommonReq
	Auth
}

func (receiver *HeartbeatReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.Auth.Format(ctx)
}

func (receiver *HeartbeatReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := receiver.Auth.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	return
}

// HeartbeatResp 心跳返回
type HeartbeatResp struct {
	PopUpResp
}

type PayLogReq struct {
	PayReq
	PayChannelId int64  `json:"pay_channel_id"`
	OrderId      string `json:"order_id"`
}

type PayLogResp struct {
	log.OdsPayLogModel
}

type PayCallbackData struct {
	CommonReq
	MerchantOrderId string
	OrderId         string
	Money           int
	Status          string
	OpenId          string
}

type PayChannelCallbackReq map[string]interface{}

func (receiver PayChannelCallbackReq) Validate() (err error) {
	if receiver.GetPayChannelId() <= 0 {
		err = errors2.New("充值渠道ID不存在")
		return
	}
	return
}

func (receiver PayChannelCallbackReq) GetPayChannelId() int64 {
	channelId, ok := receiver["pay_channel_id"]
	if !ok {
		return 0
	}
	return cast.ToInt64(channelId)
}

type PayChannelCallbackResp struct {
	Content string
}

const (
	PopUpConfigTypeRealName = "real_name" // 实名
	PopUpConfigTypePhone    = "phone"     // 绑定手机
	PopUpConfigTypeCustom   = "custom"    // 自定义
)

// PopUpConfig 弹窗配置
type PopUpConfig struct {
	Show int              `json:"show"`
	Url  string           `json:"url"`
	Btns []PopUpConfigBtn `json:"btns"`
}

const (
	PopUpConfigBtnConfirm  = "confirm"
	PopUpConfigBtnCancel   = "cancel"
	PopUpConfigBtnRealName = "real_name"
	PopUpConfigBtnPhone    = "phone"
	PopUpConfigBtnLogout   = "logout"
	PopUpConfigBtnExit     = "exit"
)

type PopUpConfigBtn struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

const (
	BuildPopUpSourceInit      = "init"
	BuildPopUpSourceLogin     = "login"
	BuildPopUpSourcePay       = "pay"
	BuildPopUpSourceHeartbeat = "heartbeat"
)

type BuildPopUpReq struct {
	CommonReq
	UserId int64
	Source string
}

type HeartbeatSwitch struct {
	Toggle   int `json:"toggle"`
	Interval int `json:"interval"`
}
