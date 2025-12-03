package media

import (
	"cngamesdk.com/api/model/api"
	"context"
	"github.com/cngamesdk/go-core/model/sql/advertising"
	"github.com/cngamesdk/go-core/validate"
	"github.com/pkg/errors"
	"strings"
)

type ReportResp struct {
	List interface{} `json:"list"`
}

type BaseReportListItemResp struct {
	Event string `json:"event"`
}

type ReportRegReq struct {
	api.CommonReq
	api.Auth
}

func (receiver *ReportRegReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.Auth.Format(ctx)
}

func (receiver *ReportRegReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	err = receiver.Auth.Validate(ctx)
	return
}

type ReportRegResp struct {
	BaseReportListItemResp
}

type ReportLoginReq struct {
	api.CommonReq
	api.Auth
}

func (receiver *ReportLoginReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.Auth.Format(ctx)
}

func (receiver *ReportLoginReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	err = receiver.Auth.Validate(ctx)
	return
}

type ReportLoginResp struct {
	BaseReportListItemResp
}

type ReportPayReq struct {
	api.CommonReq
	api.Auth
	OrderId string `json:"order_id" form:"order_id" binding:"required"`
}

func (receiver *ReportPayReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.Auth.Format(ctx)
}

func (receiver *ReportPayReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.EmptyString(receiver.OrderId); validateErr != nil {
		err = errors.Wrap(validateErr, "order_id")
		return
	}
	err = receiver.Auth.Validate(ctx)
	return
}

type ReportPayResp struct {
	BaseReportListItemResp
	Money int `json:"money"` //上报金额，单位：分
}

type ReportCallbackReq struct {
	api.CommonReq
	api.Auth
	Event   string `json:"event" form:"event" binding:"required"`
	OrderId string `json:"order_id" form:"order_id"`
	Money   int    `json:"money" form:"money"`
}

func (receiver *ReportCallbackReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.Auth.Format(ctx)
	receiver.Event = strings.TrimSpace(receiver.Event)
	receiver.OrderId = strings.TrimSpace(receiver.OrderId)
}

func (receiver *ReportCallbackReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := validate.EmptyString(receiver.Event); validateErr != nil {
		err = errors.Wrap(validateErr, "event")
		return
	}
	if receiver.Event == advertising.MediaCallbackEventPay {
		if validateErr := validate.EmptyString(receiver.OrderId); validateErr != nil {
			err = errors.Wrap(validateErr, "order_id")
			return
		}
	}
	err = receiver.Auth.Validate(ctx)
	return
}

type ReportCallbackResp struct {
}
