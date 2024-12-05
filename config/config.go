package config

import (
	"fmt"
	"github.com/shadowabi/AutoDomain_rebuild/define"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
	"github.com/spf13/viper"
)

// Config struct is a wrapper of viper
type Config struct {
	*viper.Viper
}

var GlobalConfig *viper.Viper

var C define.Configure

//// DefaultInit func is default Init Config file without config file information
//func DefaultInit() {
//	GlobalConfig = &Config{
//		viper.New(),
//	}
//
//	// Config filename (no suffix)
//	GlobalConfig.SetConfigName("config")
//
//	// Config Type
//	GlobalConfig.SetConfigType("yaml")
//
//	// Searching Path
//	GlobalConfig.AddConfigPath(".")
//	GlobalConfig.AddConfigPath("../") // For Debug
//	GlobalConfig.AddConfigPath("./config")
//
//	// Reading Config
//	err := GlobalConfig.ReadInConfig()
//	Error.HandleFatal(err)
//}

func InitConfigure(file string) {
	GlobalConfig = viper.New()
	GlobalConfig.SetConfigFile(file)
	GlobalConfig.SetConfigType("yaml")
	err := GlobalConfig.ReadInConfig()
	if err != nil {
		SaveConfig(file)
		Error.HandleFatal(fmt.Errorf("不存在 config.yaml，已生成"))
	}
	Error.HandleError(GlobalConfig.Unmarshal(&C))
}

func SaveConfig(file string) {
	GlobalConfig.Set("FofaKey", "")
	GlobalConfig.Set("QuakeKey", "")
	GlobalConfig.Set("HunterKey", "")
	GlobalConfig.Set("ZoomeyeKey", "")
	GlobalConfig.Set("PulsediveKey", "")
	GlobalConfig.Set("DaydaymapKey", "")
	GlobalConfig.SetConfigFile(file)
	err := GlobalConfig.WriteConfig()
	Error.HandleFatal(err)
}
