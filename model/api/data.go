package api

import (
	"context"
	"strings"
)

type LaunchReportReq struct {
	CommonReq
	Action string `json:"action" form:"action" binding:"required"`
}

func (receiver *LaunchReportReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.Action = strings.TrimSpace(receiver.Action)
}

func (receiver *LaunchReportReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	return
}

type LaunchReportResp struct {
}

type GameReportReq struct {
	CommonReq
	Auth
	Action     string `json:"action" form:"action" binding:"required"`
	ServerId   int64
	ServerName string
	RoleId     string
	RoleName   string
	RoleLevel  int
}

func (receiver *GameReportReq) Format(ctx context.Context) {
	receiver.CommonReq.Format(ctx)
	receiver.Auth.Format(ctx)
}

func (receiver *GameReportReq) Validate(ctx context.Context) (err error) {
	if validateErr := receiver.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	return
}

type GameReportResp struct {
}
