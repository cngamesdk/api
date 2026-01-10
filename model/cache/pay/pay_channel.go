package pay

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/model/cache"
	"cngamesdk.com/api/model/sql/pay"
	"context"
	"encoding/json"
	"github.com/duke-git/lancet/v2/cryptor"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type DimPayChannelModel struct {
	pay.DimPayChannelModel
}

func NewDimPayChannelModel() *DimPayChannelModel {
	model := &DimPayChannelModel{}
	model.DimPayChannelModel.DimPayChannelModel.Db = func() *gorm.DB {
		return global.MyDb
	}
	return model
}

// GetRateLessChannel 获取费率最少的支付渠道
func (receiver *DimPayChannelModel) GetRateLessChannel(ctx context.Context, fields string, query string, args ...interface{}) (err error) {
	cacheKey := cryptor.Md5String(cache.BuildCacheKey("dim-pay-channel-rate-less-", fields, query, args))
	cacheClient := global.MyRedis
	existsResult, existsErr := cacheClient.Exists(ctx, cacheKey).Result()
	if existsErr != nil {
		err = existsErr
		return
	}
	if existsResult > 0 {
		getResult, getErr := cacheClient.Get(ctx, cacheKey).Result()
		if getErr != nil {
			err = getErr
			return
		}
		if jsonErr := json.Unmarshal([]byte(getResult), receiver); jsonErr != nil {
			err = jsonErr
			return
		}
		return
	}
	takeErr := receiver.Db().WithContext(ctx).
		Select(fields).
		Table(receiver.TableName()).
		Where(query, args...).
		Order("rate").
		Take(receiver).Error
	if takeErr != nil {
		err = takeErr
		global.Logger.ErrorCtx(ctx, "获取异常", zap.Any("err", takeErr))
		return
	}
	jsonData, jsonErr := json.Marshal(receiver)
	if jsonErr != nil {
		err = jsonErr
		global.Logger.ErrorCtx(ctx, "JSON异常", zap.Any("err", jsonErr))
		return
	}
	if _, setErr := cacheClient.Set(ctx, cacheKey, string(jsonData), time.Minute*5).Result(); setErr != nil {
		err = setErr
		return
	}
	return
}
