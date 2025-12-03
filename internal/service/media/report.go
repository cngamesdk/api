package media

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/model/api/media"
	"cngamesdk.com/api/model/sql/advertising"
	"cngamesdk.com/api/model/sql/log"
	"context"
	advertising2 "github.com/cngamesdk/go-core/model/sql/advertising"
	log2 "github.com/cngamesdk/go-core/model/sql/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ReportService struct {
}

func (receiver *ReportService) Reg(ctx context.Context, req *media.ReportRegReq) (resp interface{}, err error) {
	var list []media.ReportRegResp
	//注册事件
	regResp := media.ReportRegResp{}
	regResp.Event = advertising2.MediaCallbackEventReg
	list = append(list, regResp)
	resp = media.ReportResp{
		List: list,
	}
	return
}

func (receiver *ReportService) Login(ctx context.Context, req *media.ReportLoginReq) (resp interface{}, err error) {
	var list []media.ReportLoginResp
	//登录事件
	loginResp := media.ReportLoginResp{}
	loginResp.Event = advertising2.MediaCallbackEventLogin
	list = append(list, loginResp)
	resp = media.ReportResp{
		List: list,
	}
	return
}

func (receiver *ReportService) Pay(ctx context.Context, req *media.ReportPayReq) (resp interface{}, err error) {
	var list []media.ReportPayResp

	model := log.NewOdsPayLogModel()
	if takeErr := model.Take(ctx, "pay_status,user_id,money", "platform_id = ? and order_id = ?", req.PlatformId, req.OrderId); takeErr != nil {
		err = takeErr
		global.Logger.ErrorCtx(ctx, "获取异常", zap.Error(takeErr))
		return
	}
	if model.UserId != req.UserId {
		err = errors.New("非法请求")
		global.Logger.ErrorCtx(ctx, "数据不一致", zap.Any("model", model), zap.Any("req", req))
		return
	}
	if model.PayStatus != log2.PayStatusSuccess {
		err = errors.New("订单未成功")
		global.Logger.ErrorCtx(ctx, "订单未成功", zap.Any("model", model))
		return
	}
	//付费事件
	payResp := media.ReportPayResp{}
	payResp.Event = advertising2.MediaCallbackEventPay
	payResp.Money = model.Money
	list = append(list, payResp)
	resp = media.ReportResp{
		List: list,
	}
	return
}

func (receiver *ReportService) Callback(ctx context.Context, req *media.ReportCallbackReq) (resp media.ReportCallbackResp, err error) {
	model := advertising.NewOdsMediaCallbackLogModel()
	model.OrderId = req.OrderId
	model.Event = req.Event
	model.Status = advertising2.MediaCallbackStatusSuccess
	model.Source = advertising2.MediaCallbackSourceSdk
	model.PlatformId = req.PlatformId
	model.UserId = req.UserId
	model.SiteId = req.SiteId
	model.GameId = req.GameId
	model.ReportMoney = req.Money
	if saveErr := model.Create(ctx); saveErr != nil {
		err = saveErr
		global.Logger.ErrorCtx(ctx, "保存异常", zap.Error(saveErr))
		return
	}
	return
}
