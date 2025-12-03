package main

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/initialization"
	"flag"
	"fmt"
	"github.com/spf13/viper"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

func main() {
	var config string
	flag.StringVar(&config, "config", "", "-config=/your/config/path")
	flag.Parse()
	if config == "" {
		panic(any("配置不能为空"))
	}
	global.ConfigPath = config
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(any(fmt.Errorf("Fatal error config file: %s \n", err)))
	}
	if err = v.Unmarshal(&global.Config); err != nil {
		panic(any(fmt.Errorf("Fatal error unmarshal file: %s \n", err)))
	}

	if initDataErr := initialization.InitConfigData(v); initDataErr != nil {
		panic(any(initDataErr))
	}

	defer global.Logger.Logger.Sync()

	initialization.Init(global.Config)
}
