package game

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/common"
	"gorm.io/gorm"
)

// DimGameModel 游戏维度
type DimGameModel struct {
	common.DimGameModel
}

func NewDimGameModel() *DimGameModel {
	model := &DimGameModel{}
	model.DimGameModel.Db = func() *gorm.DB {
		return model.DimGameModel.Db()
	}
	return model
}

func (receiver *DimGameModel) Db() *gorm.DB {
	return global.MyDb
}

func (receiver *DimGameModel) GetHashKey() string {
	return global.Config.Common.GameHashKey
}
