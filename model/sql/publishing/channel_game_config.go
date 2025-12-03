package publishing

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/publishing"
	"gorm.io/gorm"
)

// DimPublishingChannelGameConfigModel 发行渠道游戏配置表
type DimPublishingChannelGameConfigModel struct {
	publishing.DimPublishingChannelGameConfigModel
}

func NewDimPublishingChannelGameConfigModel() *DimPublishingChannelGameConfigModel {
	model := &DimPublishingChannelGameConfigModel{}
	model.DimPublishingChannelGameConfigModel.Db = func() *gorm.DB {
		return model.Db()
	}
	return model
}

func (receiver *DimPublishingChannelGameConfigModel) Db() *gorm.DB {
	return global.MyDb
}
