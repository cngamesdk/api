package game

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/common"
	"gorm.io/gorm"
)

type DimGameAppVersionConfigurationModel struct {
	common.DimGameAppVersionConfiguration
}

func NewDimGameAppVersionConfigurationModel() *DimGameAppVersionConfigurationModel {
	model := &DimGameAppVersionConfigurationModel{}
	model.DimGameAppVersionConfiguration.Db = func() *gorm.DB {
		return global.MyDb
	}
	return model
}
