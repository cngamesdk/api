package gift

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/gift"
	"gorm.io/gorm"
)

// OdsGiftClaimLogModel 礼包领取列表
type OdsGiftClaimLogModel struct {
	gift.OdsGiftClaimLogModel
}

func NewOdsGiftClaimLogModel() *OdsGiftClaimLogModel {
	model := &OdsGiftClaimLogModel{}
	model.OdsGiftClaimLogModel.Db = func() *gorm.DB {
		return model.Db()
	}
	return model
}

func (receiver *OdsGiftClaimLogModel) Db() *gorm.DB {
	return global.MyDb
}
