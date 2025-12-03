package pay

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/internal/service/pay/channel"
	"cngamesdk.com/api/model/api"
	"cngamesdk.com/api/model/sql/log"
	"context"
	log2 "github.com/cngamesdk/go-core/model/sql/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

// 记录下单日志 SavePayLog
func SavePayLog(ctx context.Context, req api.PayLogReq) (resp api.PayLogResp, err error) {
	payModel := log.NewOdsPayLogModel()
	payModel.PlatformId = req.PlatformId
	payModel.OrderId = req.OrderId
	payModel.UserId = req.UserId
	payModel.ServerId = req.ServerId
	payModel.ServerName = req.ServerName
	payModel.RoleId = req.RoleId
	payModel.RoleName = req.RoleName
	payModel.ProductId = req.ProductId
	payModel.ProductName = req.ProductName
	payModel.Money = req.Money
	payModel.PayStatus = log2.PayStatusPreorder
	payModel.PayTime = time.Now()
	payModel.Ext = req.Ext
	payModel.GameId = req.GameId
	payModel.PayChannelId = req.PayChannelId
	payModel.ChannelId = req.ChannelId
	payModel.AgentId = req.AgentId
	payModel.SiteId = req.SiteId
	payModel.ClientIp = req.ClientIp
	payModel.Ipv4 = req.Ipv4
	payModel.Ipv6 = req.Ipv6
	payModel.UserAgent = req.UserAgent
	payModel.SystemVersion = req.SystemVersion
	payModel.Brand = req.Brand
	payModel.Model = req.Model
	payModel.AndriodId = req.AndriodId
	payModel.Oaid = req.Oaid
	payModel.Imei = req.Imei
	payModel.Idfv = req.Idfv
	payModel.Network = req.Network
	payModel.AppVersionCode = req.AppVersionCode
	payModel.SdkVersionCode = req.SdkVersionCode
	if req.MediaSiteId > 0 {
		payModel.MediaSiteId = req.SiteId
		payModel.SiteId = req.MediaSiteId
	}
	if saveErr := payModel.Create(ctx); saveErr != nil {
		err = saveErr
		global.Logger.Error("保存日志异常", zap.Any("err", saveErr), zap.Any("data", payModel))
		return
	}
	resp.OdsPayLogModel = *payModel
	return
}

func GetPayChannels() []channel.PayChannelInterface {
	return []channel.PayChannelInterface{
		channel.NewWeiXinOfficialPayChannel(),
		channel.NewAlipayOfficialPayChannel(),
	}
}

// GetSdkPayChannelFactory 获取SDK支付渠道
func GetSdkPayChannelFactory(req *api.PayReq) (resp channel.PayChannelInterface, err error) {
	//根据规则切换支付渠道
	//游戏主体
	payChannels := GetPayChannels()
	payChannel := payChannels[0]
	resp = payChannel
	return
}

// PayCallback 充值回调
func PayCallback(ctx context.Context, req api.PayCallbackData) (err error) {
	if req.OrderId == "" {
		err = errors.New("订单号为空")
		return
	}
	orderModel := &log.OdsPayLogModel{}
	if takeErr := orderModel.Take(ctx, "*", "platform_id = ? and order_id = ?", req.PlatformId, req.OrderId); takeErr != nil {
		err = takeErr
		global.Logger.Error("获取订单异常", zap.Any("err", takeErr))
		return
	}
	if orderModel.Money != req.Money {
		err = errors.New("订单金额不一致")
		global.Logger.Error("订单金额不一致", zap.Any("req", req))
		return
	}
	updateModel := &log.OdsPayLogModel{}
	updateModel.PayStatus = req.Status
	updateModel.MerchantOrderId = req.MerchantOrderId
	updateModel.CallbackTime = time.Now()
	if updateErr := updateModel.Updates(ctx, "id = ?", orderModel.Id); updateErr != nil {
		err = updateErr
		global.Logger.Error("更新异常", zap.Any("err", updateErr))
		return
	}
	return
}
