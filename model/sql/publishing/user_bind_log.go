package publishing

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/publishing"
	"gorm.io/gorm"
)

// OdsPublishingUserBindLogModel 发行渠道用户绑定日志
type OdsPublishingUserBindLogModel struct {
	publishing.OdsPublishingUserBindLogModel
}

func NewOdsPublishingUserBindLogModel() *OdsPublishingUserBindLogModel {
	model := &OdsPublishingUserBindLogModel{}
	model.OdsPublishingUserBindLogModel.Db = func() *gorm.DB {
		return model.Db()
	}
	return model
}

func (receiver *OdsPublishingUserBindLogModel) Db() *gorm.DB {
	return global.MyDb
}
