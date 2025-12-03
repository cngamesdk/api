package log

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/log"
	"gorm.io/gorm"
)

type OdsRegLogModel struct {
	log.OdsRegLogModel
}

func NewOdsRegLogModel() *OdsRegLogModel {
	model := &OdsRegLogModel{}
	model.OdsRegLogModel.Db = func() *gorm.DB {
		return model.Db()
	}
	return model
}

func (receiver *OdsRegLogModel) Db() *gorm.DB {
	return global.MyDb
}
