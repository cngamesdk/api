package log

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/model/cache"
	"cngamesdk.com/api/model/sql/log"
	"context"
	"encoding/json"
	"github.com/duke-git/lancet/v2/cryptor"
	"go.uber.org/zap"
	"time"
)

type DwdGameRegLogModel struct {
	log.DwdGameRegLogModel
}

func (receiver *DwdGameRegLogModel) Take(ctx context.Context, fields string, query string, args ...interface{}) (err error) {
	cacheKey := cryptor.Md5String(cache.BuildCacheKey("dwd-game-reg-log-take-", fields, query, args))
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
	takeErr := receiver.DwdGameRegLogModel.Take(ctx, fields, query, args)
	if takeErr != nil {
		err = takeErr
		global.Logger.Error("获取异常", zap.Any("err", takeErr))
		return
	}
	jsonData, jsonErr := json.Marshal(receiver)
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
