package initialization

import (
	"cngamesdk.com/api/config"
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/middleware"
	config3 "github.com/cngamesdk/go-core/config"
	log2 "github.com/cngamesdk/go-core/log"
	translator2 "github.com/cngamesdk/go-core/translator"
	"github.com/gin-gonic/gin"
)

func Init(config config.Config) {
	global.Logger = log2.MyLogger{
		CtxRequestIdKey: global.Config.Common.CtxRequestIdKey,
	}
	global.Logger.Logger = log2.NewFileZapLogger(config.Log)

	//初始化数据库
	db, dbErr := config3.OpenMysql(config.Mysql)

	if dbErr != nil {
		panic(any(dbErr))
	}
	global.MyDb = db

	if migrateErr := Migrate(); migrateErr != nil {
		panic(any(migrateErr))
	}

	//初始化REDIS
	myRedis, myRedisErr := config3.OpenRedis(config.Redis)
	if myRedisErr != nil {
		panic(any(myRedisErr))
	}
	global.MyRedis = myRedis

	//初始化服务
	r := gin.New()

	r.Use(middleware.Trace(), middleware.Recovery())

	//路由初始化
	RouteInit(r)

	//初始化翻译器
	if err := translator2.InitTrans("zh"); err != nil {
		panic(any(err))
	}

	if err := r.Run(config.Server.Host); err != nil {
		panic(any(err))
	}
}
