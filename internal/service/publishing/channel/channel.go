package channel

import "github.com/cngamesdk/go-core/model/sql"

type baseChannel struct {
	config sql.JSON
}

func (receiver baseChannel) SetConfig(req sql.JSON) {
	receiver.config = req
}

func (receiver baseChannel) GetConfig() (resp sql.JSON) {
	return receiver.config
}
