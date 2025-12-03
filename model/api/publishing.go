package api

import (
	"context"
	"github.com/cngamesdk/go-core/validate"
	errors2 "github.com/pkg/errors"
	"github.com/spf13/cast"
	"strings"
)

type PublishingLoginReq struct {
	CommonReq
	OpenId string `json:"open_id" form:"open_id" binding:"required"`
	Token  string `json:"token" form:"token" binding:"required"`
	Ext    string `json:"ext" form:"ext"`
}

func (receiver *PublishingLoginReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.OpenId = strings.TrimSpace(receiver.OpenId)
	receiver.Token = strings.TrimSpace(receiver.Token)
	receiver.Ext = strings.TrimSpace(receiver.Ext)
}

func (receiver *PublishingLoginReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if receiver.GameId <= 0 {
		err = errors2.New("发行游戏ID为空")
		return
	}
	if receiver.ChannelId <= 0 {
		err = errors2.New("发行ID为空")
		return
	}
	if receiver.AgentId <= 0 {
		err = errors2.New("渠道ID为空")
		return
	}
	if receiver.SiteId <= 0 {
		err = errors2.New("广告位ID为空")
		return
	}
	if validateErr := validate.EmptyString(receiver.OpenId); validateErr != nil {
		err = errors2.Wrap(validateErr, "open_id")
		return
	}
	if validateErr := validate.EmptyString(receiver.Token); validateErr != nil {
		err = errors2.Wrap(validateErr, "token")
		return
	}
	return
}

type PublishingLoginResp struct {
	PopUpResp
	BaseUserAuthRespModel
}

type PublishingPayReq struct {
	PayReq
}

func (receiver PublishingPayReq) Format(ctx context.Context) {
	receiver.PayReq.Format(ctx)
}

func (receiver PublishingPayReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if receiver.GameId <= 0 {
		err = errors2.New("发行游戏ID为空")
		return
	}
	if receiver.ChannelId <= 0 {
		err = errors2.New("发行ID为空")
		return
	}
	if receiver.AgentId <= 0 {
		err = errors2.New("渠道ID为空")
		return
	}
	if receiver.SiteId <= 0 {
		err = errors2.New("广告位ID为空")
		return
	}
	if validateErr := receiver.PayReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	return
}

type PublishingPayResp struct {
	OrderId string `json:"order_id"`
	Ext     string `json:"ext"`
	PopUpResp
}

type PublishingChannelPayCallbackReq map[string]interface{}

func (receiver PublishingChannelPayCallbackReq) Validate() (err error) {
	if receiver.GetChannelId() <= 0 {
		err = errors2.New("渠道ID为空")
		return
	}
	return
}

func (receiver PublishingChannelPayCallbackReq) GetChannelId() int64 {
	payChannelId, ok := receiver["channel_id"]
	if !ok {
		return 0
	}
	return cast.ToInt64(payChannelId)
}

type PublishingChanelPayCallbackResp struct {
	Content string
}
