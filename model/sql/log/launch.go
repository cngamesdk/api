package log

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/log"
	"gorm.io/gorm"
)

type OdsLaunchLogModel struct {
	log.OdsLaunchLogModel
}

func NewOdsLaunchLogModel() *OdsLaunchLogModel {
	model := &OdsLaunchLogModel{}
	model.OdsLaunchLogModel.Db = func() *gorm.DB {
		return model.Db()
	}
	return model
}

func (receiver *OdsLaunchLogModel) Db() *gorm.DB {
	return global.MyDb
}
