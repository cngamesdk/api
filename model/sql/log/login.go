package log

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/log"
	"gorm.io/gorm"
)

type OdsLoginLogModel struct {
	log.OdsLoginLogModel
}

func NewOdsLoginLogModel() *OdsLoginLogModel {
	model := &OdsLoginLogModel{}
	model.OdsLoginLogModel.Db = func() *gorm.DB {
		return model.Db()
	}
	return model
}

func (receiver *OdsLoginLogModel) Db() *gorm.DB {
	return global.MyDb
}
