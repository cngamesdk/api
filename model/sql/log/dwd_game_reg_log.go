package log

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/log"
	"gorm.io/gorm"
)

type DwdGameRegLogModel struct {
	log.DwdGameRegLogModel
}

func NewDwdGameRegLogModel() *DwdGameRegLogModel {
	model := &DwdGameRegLogModel{}
	model.DwdGameRegLogModel.Db = func() *gorm.DB {
		return global.MyDb
	}
	return model
}
