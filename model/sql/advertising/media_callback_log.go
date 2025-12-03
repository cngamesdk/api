package advertising

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/advertising"
	"gorm.io/gorm"
)

type OdsMediaCallbackLogModel struct {
	advertising.OdsMediaCallbackLogModel
}

func NewOdsMediaCallbackLogModel() *OdsMediaCallbackLogModel {
	model := &OdsMediaCallbackLogModel{}
	model.OdsMediaCallbackLogModel.Db = func() *gorm.DB {
		return global.MyDb
	}
	return model
}
