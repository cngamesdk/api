package sms

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/model/api"
	"cngamesdk.com/api/model/sql/user"
	"context"
	"github.com/cngamesdk/go-core/goroutine"
	"github.com/cngamesdk/go-core/sms"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/duke-git/lancet/v2/random"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"time"
)

type SmsService struct {
	Phone string
}

func NewSmsService(phone string) *SmsService {
	return &SmsService{
		Phone: phone,
	}
}

func (receiver *SmsService) getCacheKey() string {
	return cryptor.Md5String(global.Config.Common.CommonHashKey + receiver.Phone)
}

func (receiver *SmsService) SendCheckCode(ctx context.Context, req *api.SmsSendReq) (err error) {
	receiver.Phone = req.Phone
	checkCode := random.RandNumeral(6)
	if slices.Contains(global.Config.Common.TestPhones, receiver.Phone) {
		checkCode = "123456"
	}
	smsService := sms.NewSmsService(receiver.Phone)
	sendErr := smsService.SendCheckCode(ctx, checkCode)
	if sendErr != nil {
		err = sendErr
		global.Logger.Error("发送验证码异常", zap.Any("err", sendErr), zap.Any("data", receiver.Phone))
		return
	}
	cacheKey := receiver.getCacheKey()
	_, cacheSetResultErr := global.MyRedis.Set(ctx, cacheKey, checkCode, time.Minute*5).Result()
	if cacheSetResultErr != nil {
		err = cacheSetResultErr
		global.Logger.Error("设置缓存异常", zap.Any("err", cacheSetResultErr))
		return
	}

	goroutine.CreateGoroutine(func() {

		sqlModel := user.NewOdsUserSmsSendLogModel()
		sqlModel.PlatformId = req.PlatformId
		sqlModel.Phone = req.Phone
		sqlModel.ActionTime = time.Now()
		sqlModel.Content = smsService.GetContent()
		sqlModel.Result = smsService.GetResult()

		if recordErr := sqlModel.Create(ctx); recordErr != nil {
			global.Logger.Error("记录日志异常", zap.Error(recordErr), zap.Any("data", sqlModel))
		}
	}, func(any2 any) {
		global.Logger.Error("协程异常", zap.Any("err", any2))
	})

	return
}

func (receiver SmsService) VerifyCheckCode(ctx context.Context, code string) (err error) {
	cacheKey := receiver.getCacheKey()
	existsResult, existsErr := global.MyRedis.Exists(ctx, cacheKey).Result()
	if existsErr != nil {
		err = existsErr
		global.Logger.Error("EXISTS缓存异常", zap.Any("err", existsErr))
		return
	}
	if existsResult <= 0 {
		err = errors.New("请先发送验证码")
		return
	}
	getResult, getResultErr := global.MyRedis.Get(ctx, cacheKey).Result()
	if getResultErr != nil {
		err = getResultErr
		global.Logger.Error("读取缓存异常", zap.Any("err", getResultErr))
		return
	}
	if getResult != code {
		err = errors.New("验证码错误")
		global.Logger.Warn("验证码验证错误", zap.Any("cache", getResult), zap.Any("code", code))
		return
	}
	global.MyRedis.Del(ctx, cacheKey)
	return
}
