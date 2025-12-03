package api

import (
	"context"
	"github.com/cngamesdk/go-core/validate"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"math"
	"time"
)

// 研发二次验证请求参数
type CpAuthReq struct {
	CommonReq
	UserId    int64  `json:"user_id" form:"user_id" binding:"required"`
	Token     string `json:"token" form:"token" binding:"required"`
	Timestamp int64  `json:"timestamp" form:"timestamp" binding:"required"`
	Sign      string `json:"sign" form:"sign" binding:"required"`
}

func (c *CpAuthReq) Format(ctx context.Context) {
	c.CommonReq.Format(ctx)
}

func (c *CpAuthReq) Validate(ctx context.Context) (err error) {
	if validateErr := c.CommonReq.Validate(ctx); validateErr != nil {
		err = validateErr
		return
	}
	if c.PlatformId <= 0 {
		err = errors.New("平台ID为空")
		return
	}
	if c.GameId <= 0 {
		err = errors.New("游戏ID为空")
		return
	}
	if c.UserId <= 0 {
		err = errors.New("用户ID为空")
		return
	}
	if validateErr := validate.EmptyString(c.Token); validateErr != nil {
		err = validateErr
		return
	}
	if c.Timestamp <= 0 {
		err = errors.New("时间戳为空")
		return
	}
	if math.Abs(cast.ToFloat64(time.Now().Unix()-c.Timestamp)) > 300 {
		err = errors.New("请求已经过期")
		return
	}
	if validateErr := validate.EmptyString(c.Sign); validateErr != nil {
		err = validateErr
		return
	}
	return
}

type CpAuthResp struct {
	UserId int64 `json:"user_id" form:"user_id" binding:"required"`
}
