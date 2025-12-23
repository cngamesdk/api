package config

import (
	"github.com/cngamesdk/go-core/config"
)

type Config struct {
	Installed  int                 `mapstructure:"installed" json:"installed" yaml:"installed"`
	Common     config.CommonConfig `mapstructure:"common" json:"common" yaml:"common"`
	Mysql      config.MySql        `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis      config.Redis        `mapstructure:"redis" json:"redis" yaml:"redis"`
	Log        config.FileLog      `mapstructure:"log" json:"log" yaml:"log"`
	Server     config.Server       `mapstructure:"server" json:"server" yaml:"server"`
	DataReport DataReport          `mapstructure:"data_report" json:"data_report" yaml:"data_report"`
}

type DataReport struct {
	RecordLaunchLog int `mapstructure:"record_launch_log" json:"record_launch_log" yaml:"record_launch_log"` // 是否记录启动日志
}
