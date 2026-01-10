package pay

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/common"
	"gorm.io/gorm"
)

type DimPayChannelModel struct {
	common.DimPayChannelModel
}

func NewDimPayChannelModel() *DimPayChannelModel {
	model := &DimPayChannelModel{}
	model.DimPayChannelModel.Db = func() *gorm.DB {
		return global.MyDb
	}
	return model
}
