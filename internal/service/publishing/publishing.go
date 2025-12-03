package publishing

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/internal/service/pay"
	channel2 "cngamesdk.com/api/internal/service/pay/channel"
	"cngamesdk.com/api/internal/service/publishing/channel"
	"cngamesdk.com/api/internal/service/third"
	"cngamesdk.com/api/model/api"
	"cngamesdk.com/api/model/sql/publishing"
	user2 "cngamesdk.com/api/model/sql/user"
	"context"
	"encoding/json"
	"github.com/cngamesdk/go-core/model/sql"
	"github.com/cngamesdk/go-core/validate"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type ChannelInterface interface {
	GetId() int64
	SetConfig(req sql.JSON)
	GetConfig() (resp sql.JSON)
	Login(ctx context.Context, req *api.PublishingLoginReq) (resp string, err error)
	Pay(ctx context.Context, req *api.PublishingPayReq) (resp string, err error)
	PayCallback(ctx context.Context, req api.PublishingChannelPayCallbackReq) (resp api.PayCallbackData, err error)
}

// GetChannel 获取联运渠道
func GetChannel(channelId int64) ChannelInterface {
	testChannel := channel.TestChannel{}
	if testChannel.GetId() == channelId {
		return testChannel
	}
	return nil
}

type PublishingService struct {
}

// Login 联运渠道登录
func (receiver PublishingService) Login(ctx context.Context, req *api.PublishingLoginReq) (resp *user2.OdsUserInfoLogModel, err error) {
	configModel := publishing.NewDimPublishingChannelGameConfigModel()
	if takeErr := configModel.Take(ctx, "*", "platform_id = ? and game_id = ? and channel_id = ? and site_id = ?",
		req.PlatformId, req.GameId, req.ChannelId, req.SiteId); takeErr != nil {
		err = takeErr
		global.Logger.ErrorCtx(ctx, "获取异常", zap.Any("err", takeErr))
		return
	}
	channelService := GetChannel(configModel.ChannelId)
	if channelService == nil {
		err = errors.New("发行渠道未配置")
		global.Logger.ErrorCtx(ctx, "发行渠道未配置", zap.Any("req", req))
		return
	}
	channelService.SetConfig(configModel.Config)
	openId, loginErr := channelService.Login(ctx, req)
	if loginErr != nil {
		err = errors.New("远程登录失败")
		global.Logger.ErrorCtx(ctx, "远程登录失败", zap.Any("err", loginErr))
		return
	}
	if req.OpenId != openId {
		err = errors.New("open_id不一致")
		global.Logger.ErrorCtx(ctx, "open_id不一致", zap.Any("data", openId))
		return
	}

	thirdService := &third.ThirdService{}
	authResult, thirdAuthErr := thirdService.Login(ctx, req)
	if thirdAuthErr != nil {
		err = thirdAuthErr
		global.Logger.ErrorCtx(ctx, "授权异常", zap.Any("err", thirdAuthErr))
		return
	}
	resp = authResult
	return
}

// Pay 联运渠道支付
func (receiver PublishingService) Pay(ctx context.Context, req *api.PublishingPayReq) (resp api.PublishingPayResp, err error) {
	orderId := (channel2.PayChannel{}).CreateOrderId()
	payLogReq := api.PayLogReq{}
	payLogReq.PayReq = req.PayReq
	payLogReq.PayChannelId = global.Config.Common.PublishingPayChannelId
	payLogReq.OrderId = orderId
	saveResult, saveErr := pay.SavePayLog(ctx, payLogReq)
	if saveErr != nil {
		err = saveErr
		global.Logger.Error("保存日志异常", zap.Any("err", saveErr), zap.Any("data", payLogReq))
		return
	}

	extData := ExtData{}
	extData.Id = saveResult.Id
	extData.Sign = cryptor.Md5String(global.Config.Common.CommonHashKey + cast.ToString(saveResult.Id))
	buildExt, buildErr := BuildPayExt(extData)
	if buildErr != nil {
		err = buildErr
		global.Logger.ErrorCtx(ctx, "构建EXT异常", zap.Error(buildErr))
	}
	resp.OrderId = orderId
	resp.Ext = buildExt
	return
}

type ExtData struct {
	Id   int64  `json:"id"`
	Sign string `json:"sign"`
}

// BuildPayExt 构建支付的Ext
func BuildPayExt(req ExtData) (resp string, err error) {
	jsonData, jsonErr := json.Marshal(req)
	if jsonErr != nil {
		err = jsonErr
		return
	}
	resp = cryptor.Base64StdEncode(string(cryptor.AesEcbEncrypt(jsonData, []byte(global.Config.Common.AesCryptKey))))
	return
}

// ParsePayExt 解析支付的Ext
func ParsePayExt(req string) (resp ExtData, err error) {
	if validateErr := validate.EmptyString(req); validateErr != nil {
		err = errors.Wrap(validateErr, "ext")
		return
	}
	decryptData := cryptor.AesEcbDecrypt([]byte(cryptor.Base64StdDecode(req)), []byte(global.Config.Common.AesCryptKey))
	if validateErr := validate.EmptyString(string(decryptData)); validateErr != nil {
		err = errors.Wrap(validateErr, "ext解密失败")
		return
	}
	jsonErr := json.Unmarshal(decryptData, &resp)
	if jsonErr != nil {
		err = jsonErr
		return
	}
	mySign := cryptor.Md5String(global.Config.Common.CommonHashKey + cast.ToString(resp.Id))
	if mySign != resp.Sign {
		err = errors.New("sign校验异常")
		global.Logger.Error("sign校验异常", zap.Any("data", resp), zap.Any("mySign", mySign))
		return
	}
	return
}
