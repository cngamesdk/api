package user

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/model/sql/user"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type OdsUserInfoLogModel struct {
	user.OdsUserInfoLogModel
}

func NewOdsUserInfoLogModel() *OdsUserInfoLogModel {
	model := &OdsUserInfoLogModel{}
	model.OdsUserInfoLogModel = *user.NewOdsUserInfoLogModel()
	return model
}

func (receiver *OdsUserInfoLogModel) getCacheKey(userId int64) string {
	return fmt.Sprintf("user-info-%d", userId)
}

func (receiver *OdsUserInfoLogModel) DelById(ctx context.Context, userId int64) {
	cacheKey := receiver.getCacheKey(userId)
	cacheClient := global.MyRedis
	_, actionErr := cacheClient.Del(ctx, cacheKey).Result()
	if actionErr != nil {
		global.Logger.Error("删除异常", zap.Any("err", actionErr))
		return
	}
	return
}

func (receiver *OdsUserInfoLogModel) TakeById(ctx context.Context, userId int64) (err error) {
	cacheKey := receiver.getCacheKey(userId)
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
	takeErr := receiver.OdsUserInfoLogModel.OdsUserInfoLogModel.Take(ctx, "*", "id = ?", userId)
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
