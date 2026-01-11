package pay

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/model/cache"
	"cngamesdk.com/api/model/sql/pay"
	"context"
	"encoding/json"
	error2 "github.com/cngamesdk/go-core/model/error"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type DimPayChannelSwitchModel struct {
	pay.DimPayChannelSwitchModel
}

func NewDimPayChannelSwitchModel() *DimPayChannelSwitchModel {
	model := &DimPayChannelSwitchModel{}
	model.DimPayChannelSwitchModel.DimPayChannelSwitchModel.Db = func() *gorm.DB {
		return global.MyDb
	}
	return model
}

func (receiver *DimPayChannelSwitchModel) Db() *gorm.DB {
	return global.MyDb
}

func (receiver *DimPayChannelSwitchModel) FindAllRules(ctx context.Context, fields string, query string, args ...interface{}) (
	resp []DimPayChannelSwitchModel, err error) {
	cacheKey := cryptor.Md5String(cache.BuildCacheKey("dim-pay-channel-switch-", fields, query, args))
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
		if jsonErr := json.Unmarshal([]byte(getResult), &resp); jsonErr != nil {
			err = jsonErr
			return
		}
		return
	}
	takeErr := receiver.Db().WithContext(ctx).
		Select(fields).
		Table(receiver.TableName()).
		Where(query, args...).
		Order("sort DESC, id DESC").
		Find(&resp).Error
	if takeErr != nil {
		if !errors.Is(takeErr, error2.ErrorRecordIsNotFind) {
			err = takeErr
			global.Logger.Error("获取异常", zap.Any("err", takeErr))
			return
		}
	}
	jsonData, jsonErr := json.Marshal(resp)
	if jsonErr != nil {
		err = jsonErr
		global.Logger.Error("JSON异常", zap.Any("err", jsonErr))
		return
	}
	if _, setErr := cacheClient.Set(ctx, cacheKey, string(jsonData), time.Minute*5).Result(); setErr != nil {
		err = setErr
		return
	}
	return
}
