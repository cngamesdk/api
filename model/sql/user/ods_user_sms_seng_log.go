package user

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/user"
	"gorm.io/gorm"
)

// OdsUserSmsSendLogModel 用户发送短信日志
type OdsUserSmsSendLogModel struct {
	user.OdsUserSmsSendLogModel
}

func NewOdsUserSmsSendLogModel() *OdsUserSmsSendLogModel {
	model := &OdsUserSmsSendLogModel{}
	model.OdsUserSmsSendLogModel.Db = func() *gorm.DB {
		return model.Db()
	}
	return model
}

func (receiver *OdsUserSmsSendLogModel) Db() *gorm.DB {
	return global.MyDb
}
