package cp

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/internal/service/token"
	"cngamesdk.com/api/model/api"
	game3 "cngamesdk.com/api/model/cache/game"
	"cngamesdk.com/api/model/sql/log"
	"context"
	"fmt"
	"github.com/cngamesdk/go-core/model/sql/common"
	log2 "github.com/cngamesdk/go-core/model/sql/log"
	"github.com/cngamesdk/go-core/validate"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/duke-git/lancet/v2/netutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"
	"time"
)

type CpService struct {
}

func (receiver CpService) PayCallback(ctx context.Context, req *api.PayCallbackData) (err error) {
	orderModel := log.NewOdsPayLogModel()
	if takeErr := orderModel.Take(ctx, "*", "platform_id = ? and order_id = ?", req.PlatformId, req.OrderId); takeErr != nil {
		err = takeErr
		global.Logger.Error("获取异常", zap.Any("err", takeErr))
		return
	}
	if orderModel.PayStatus != log2.PayStatusSuccess {
		err = errors.New("订单状态未成功")
		global.Logger.Warn("订单未成功", zap.Any("data", orderModel))
		return
	}
	if orderModel.CpStatus == log2.CpStatusSuccess {
		return
	}

	gameRate := 0
	cpCallbackUrl := ""
	cpCallbackResult := ""

	defer func() {
		updateModel := &log.OdsPayLogModel{}
		updateModel.CpUrl = cpCallbackUrl
		updateModel.CpStatus = log2.CpStatusSuccess
		if err != nil {
			updateModel.CpStatus = log2.CpStatusFail
		}
		updateModel.CpResult = cpCallbackResult
		updateModel.CpSendTime = time.Now()
		updateModel.CpSendRepeatTimes = orderModel.CpSendRepeatTimes + 1
		updateModel.CpMoney = orderModel.Money * gameRate
		if updateErr := updateModel.Updates(ctx, "id = ?", orderModel.Id); updateErr != nil {
			err = updateErr
			global.Logger.Error("更新异常", zap.Any("err", updateErr))
			return
		}
		return
	}()

	gameModel := game3.NewDimGameModel()
	if takeErr := gameModel.Take(ctx, "*", "id = ?", orderModel.GameId); takeErr != nil {
		err = takeErr
		global.Logger.ErrorCtx(ctx, "获取游戏信息异常", zap.Error(takeErr))
		return
	}
	gameRate = gameModel.GameRate

	cpGameId := orderModel.GameId

	if *gameModel.CpGameId > 0 {
		cpGameId = *gameModel.CpGameId
	}
	gameProductConfigModel := game3.NewDimGameAppVersionConfigurationWithProductConfigModel()
	if takeErr := gameProductConfigModel.Take(ctx, "*", "platform_id = ? and game_id = ?", orderModel.PlatformId, orderModel.GameId); takeErr != nil {
		err = takeErr
		global.Logger.ErrorCtx(ctx, "获取配置异常", zap.Error(takeErr))
		return
	}
	cpCallbackUrl = gameProductConfigModel.ProductConfig.ShippingAddress
	if validateErr := validate.EmptyString(cpCallbackUrl); validateErr != nil {
		err = errors.Wrap(validateErr, "回调地址")
		global.Logger.ErrorCtx(ctx, "回调地址", zap.Error(err))
		return
	}
	if strings.Index(cpCallbackUrl, "?") > 0 {
		cpCallbackUrl = cpCallbackUrl + "&"
	} else {
		cpCallbackUrl = cpCallbackUrl + "?"
	}
	timestamp := time.Now().Unix()
	signKey := common.GetGamePayKey(cpGameId, global.Config.Common.GameHashKey)
	sign := fmt.Sprintf("%d%s%d%s%s%d%d%s",
		orderModel.UserId, orderModel.OrderId, orderModel.ServerId, orderModel.RoleId, orderModel.Ext, orderModel.Money, timestamp, signKey)
	buildMap := map[string]interface{}{
		"user_id":   orderModel.UserId,
		"order_id":  orderModel.OrderId,
		"server_id": orderModel.ServerId,
		"role_id":   orderModel.RoleId,
		"ext":       orderModel.Ext,
		"money":     orderModel.Money,
		"timestamp": timestamp,
		"sign":      sign,
	}
	cpCallbackUrl = cpCallbackUrl + netutil.ConvertMapToQueryString(buildMap)

	httpClientConfig := &netutil.HttpClientConfig{
		Timeout:          time.Second * 5,
		HandshakeTimeout: time.Second * 5,
		ResponseTimeout:  time.Second * 5,
	}
	httpRequest := &netutil.HttpRequest{}
	httpRequest.Method = "GET"
	httpRequest.RawURL = cpCallbackUrl
	httpClient := netutil.NewHttpClientWithConfig(httpClientConfig)
	httpClient.Context = ctx

	monitorStartTime := time.Now()
	httpResponse, httpErr := httpClient.SendRequest(httpRequest)
	global.Logger.Info("发货记录", zap.Any("url", cpCallbackUrl), zap.Any("耗时", time.Now().Sub(monitorStartTime).Seconds()))

	if httpErr != nil {
		err = httpErr
		global.Logger.Error("发货异常", zap.Any("err", httpErr))
		return
	}
	var responseByte []byte
	_, readErr := httpResponse.Body.Read(responseByte)
	if readErr != nil {
		err = readErr
		global.Logger.Error("读取异常", zap.Any("err", readErr))
		return
	}
	cpCallbackResult = string(responseByte)

	type responseStruct struct {
		code int
		msg  string
	}
	var myResponse responseStruct
	if decodeErr := httpClient.DecodeResponse(httpResponse, &myResponse); decodeErr != nil {
		err = decodeErr
		global.Logger.Error("解码异常", zap.Any("err", decodeErr))
		return
	}
	if myResponse.code != 1 {
		err = errors.New("接口返回：" + myResponse.msg)
		global.Logger.Error("解码异常", zap.Any("data", myResponse))
		return
	}
	return
}

// Auth 授权登录验证
func (receiver CpService) Auth(ctx context.Context, req *api.CpAuthReq) (resp api.CpAuthResp, err error) {
	loginKey := common.GetGameLoginKey(req.GameId, global.Config.Common.GameHashKey)
	signStr := fmt.Sprintf("%s%d%d%s%d%s", req.Platform, req.GameId, req.UserId, req.Token, req.Timestamp, loginKey)
	println(signStr)
	mySign := cryptor.Md5String(signStr)
	if req.Sign != mySign {
		err = errors.New(fmt.Sprintf("加密异常.ori: %s, dst: %s", req.Sign, mySign))
		global.Logger.Warn("加密校验异常", zap.Any("data", fmt.Sprintf("sign:%s,mySign:%s,str:%s", req.Sign, mySign, signStr)))
		return
	}
	tokenService := &token.TokenService{}
	_, tokenErr := tokenService.Verify(ctx, req.Token)
	if tokenErr != nil {
		err = tokenErr
		global.Logger.Error("验证异常", zap.Any("err", tokenErr))
		return
	}
	resp.UserId = tokenService.UserId
	return
}
