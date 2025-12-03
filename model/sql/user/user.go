package user

import (
	"cngamesdk.com/api/global"
	"context"
	"github.com/cngamesdk/go-core/goroutine"
	"github.com/cngamesdk/go-core/model/sql/user"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OdsUserInfoLogModel struct {
	user.OdsUserInfoLogModel
}

func NewOdsUserInfoLogModel() *OdsUserInfoLogModel {
	model := &OdsUserInfoLogModel{}
	model.OdsUserInfoLogModel.Db = func() *gorm.DB {
		return model.Db()
	}
	model.OdsUserInfoLogModel.GetHashKey = func() string {
		return model.GetHashKey()
	}
	model.OdsUserInfoLogModel.GetAesKey = func() string {
		return model.GetAesKey()
	}
	model.OdsUserInfoLogModel.UpdateHook = func(tx *gorm.DB) (err error) {
		goroutine.CreateGoroutine(func() {
			ctx := context.Background()
			if saveErr := model.OdsUserInfoLogModel.SaveUserOperationLog(ctx); saveErr != nil {
				global.Logger.ErrorCtx(ctx, "保存日志异常", zap.Error(saveErr))
			}
		}, func(any2 any) {
			global.Logger.Error("协程异常", zap.Any("err", any2))
		})
		return
	}
	return model
}

func (receiver *OdsUserInfoLogModel) Db() *gorm.DB {
	return global.MyDb
}

func (receiver *OdsUserInfoLogModel) GetAesKey() string {
	return global.Config.Common.AesCryptKey
}

func (receiver *OdsUserInfoLogModel) GetHashKey() string {
	return global.Config.Common.SqlMd5CryptKey
}
