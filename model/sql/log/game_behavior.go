package log

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/log"
	"gorm.io/gorm"
)

type OdsGameBehaviorLogModel struct {
	log.OdsGameBehaviorLogModel
}

func NewOdsGameBehaviorLogModel() *OdsGameBehaviorLogModel {
	model := &OdsGameBehaviorLogModel{}
	model.OdsGameBehaviorLogModel.Db = func() *gorm.DB {
		return model.Db()
	}
	return model
}

func (receiver *OdsGameBehaviorLogModel) Db() *gorm.DB {
	return global.MyDb
}
