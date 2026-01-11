package game

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/model/cache"
	"cngamesdk.com/api/model/sql/game"
	"context"
	"encoding/json"
	"github.com/cngamesdk/go-core/model/sql/common"
	"github.com/duke-git/lancet/v2/cryptor"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type DimGameModel struct {
	game.DimGameModel
	MainGameId int `json:"main_game_id" gorm:"column:main_game_id;default:0;comment:主游戏ID"`
	RootGameId int `json:"root_game_id" gorm:"column:root_game_id;default:0;comment:根游戏ID"`
}

func NewDimGameModel() *DimGameModel {
	model := &DimGameModel{}
	model.DimGameModel.DimGameModel.Db = func() *gorm.DB {
		return model.DimGameModel.Db()
	}
	return model
}

func (receiver *DimGameModel) Take(ctx context.Context, fields string, query string, args ...interface{}) (err error) {
	cacheKey := cryptor.Md5String(cache.BuildCacheKey("dim-game-take-", fields, query, args))
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
	takeErr := receiver.DimGameModel.Take(ctx, fields, query, args)
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

// DetailInfo 获取子游戏详情
func (receiver *DimGameModel) DetailInfoByGameId(ctx context.Context, gameId int64) (err error) {
	cacheKey := cryptor.Md5String(cache.BuildCacheKey("dim-game-detail-info-", gameId))
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
	takeErr := receiver.Db().WithContext(ctx).Select("game.*, main_game.root_game_id").
		Table(receiver.TableName()+" as game").Joins((&common.DimMainGameModel{}).TableName()+" as main_game ON game.platform_id = main_game.platform_id AND game.main_id = main_game.id").
		Where("game.id = ?", gameId).Take(receiver).Error
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
