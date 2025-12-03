package user

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/user"
	"gorm.io/gorm"
)

// OdsUserOperationLogModel 用户操作日志
type OdsUserOperationLogModel struct {
	user.OdsUserOperationLogModel
}

func NewOdsUserOperationLogModel() *OdsUserOperationLogModel {
	model := &OdsUserOperationLogModel{}
	model.OdsUserOperationLogModel.Db = func() *gorm.DB {
		return model.Db()
	}
	return model
}

func (receiver *OdsUserOperationLogModel) Db() *gorm.DB {
	return global.MyDb
}
