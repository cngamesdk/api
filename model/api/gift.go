package api

import (
	"context"
	errors2 "github.com/pkg/errors"
)

type GiftListReq struct {
	CommonReq
	Page     int
	PageSize int
	Auth
}

func (receiver *GiftListReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	if receiver.Page <= 0 {
		receiver.Page = 1
	}
	if receiver.PageSize <= 0 {
		receiver.PageSize = 10
	}
	receiver.Auth.Format(ctx)
}

func (receiver *GiftListReq) Validate(ctx context.Context) (err error) {
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

type GiftListResp struct {
	Page        int         `json:"page"`
	PageSize    int         `json:"page_size"`
	TotalRecord int         `json:"total_record"`
	TotalPage   int         `json:"total_page"`
	List        interface{} `json:"list"`
}

type GiftListDataResp struct {
	Id     int64  `json:"id"`
	Icon   string `json:"icon"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Status string `json:"status"`
	Code   string `json:"code"`
}

type GiftClaimReq struct {
	CommonReq
	GiftId int64 `json:"gift_id" form:"gift_id" binding:"required"`
	Auth
}

func (receiver *GiftClaimReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.Auth.Format(ctx)
}

func (receiver *GiftClaimReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if validateErr := receiver.Auth.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if receiver.GiftId <= 0 {
		err = errors2.New("礼包ID为空")
		return
	}
	return
}

type GiftClaimResp struct {
	Code string `json:"code"`
}
