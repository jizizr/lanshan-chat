package initialize

import (
	"github.com/spf13/viper"
	"lanshan_chat/app/api/global"
)

func SetupViper() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.SetConfigFile("./manifest/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic("viper read config failed, err: " + err.Error())
	}
	err = viper.Unmarshal(&global.Config)
	if err != nil {
		panic("viper unmarshal config failed, err: " + err.Error())
	}
}
