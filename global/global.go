package global

import (
	"cngamesdk.com/api/config"
	log2 "github.com/cngamesdk/go-core/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	CtxKeyClient = "ctx-client" //上下文中的客户端
)

var (
	ConfigPath string
	Logger     log2.MyLogger
	Config     config.Config
	MyDb       *gorm.DB
	MyRedis    *redis.Client
)
