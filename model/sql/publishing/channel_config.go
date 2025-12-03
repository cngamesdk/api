package publishing

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/publishing"
	"gorm.io/gorm"
)

// DimPublishingChannelConfigModel 发行渠道配置表
type DimPublishingChannelConfigModel struct {
	publishing.DimPublishingChannelConfigModel
}

func NewDimPublishingChannelConfigModel() *DimPublishingChannelConfigModel {
	model := &DimPublishingChannelConfigModel{}
	model.DimPublishingChannelConfigModel.Db = func() *gorm.DB {
		return model.Db()
	}
	return model
}

func (receiver *DimPublishingChannelConfigModel) Db() *gorm.DB {
	return global.MyDb
}
