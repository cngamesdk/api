package game

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/common"
	"gorm.io/gorm"
)

type DimProductCommonConfigurationModel struct {
	common.DimProductCommonConfigurationModel
}

func NewDimProductCommonConfigurationModel() *DimProductCommonConfigurationModel {
	model := &DimProductCommonConfigurationModel{}
	model.DimProductCommonConfigurationModel.Db = func() *gorm.DB {
		return global.MyDb
	}
	return model
}
