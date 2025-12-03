package log

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/log"
	"gorm.io/gorm"
)

type OdsPayLogModel struct {
	log.OdsPayLogModel
}

func NewOdsPayLogModel() *OdsPayLogModel {
	model := &OdsPayLogModel{}
	model.OdsPayLogModel.Db = func() *gorm.DB {
		return model.Db()
	}
	return model
}

func (receiver *OdsPayLogModel) Db() *gorm.DB {
	return global.MyDb
}
