package pay

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/common"
	"gorm.io/gorm"
)

type DimPayChannelSwitchModel struct {
	common.DimPayChannelSwitchModel
}

func NewDimPayChannelSwitchModel() *DimPayChannelSwitchModel {
	model := &DimPayChannelSwitchModel{}
	model.DimPayChannelSwitchModel.Db = func() *gorm.DB {
		return global.MyDb
	}
	return model
}
