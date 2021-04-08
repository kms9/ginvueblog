package setup

import (
	"ginvueblog/utils"
	"github.com/spf13/viper"
)

var ReConfig  *viper.Viper

// StartLogger 初始日志
func StartConfig() error {

	endpoint := "http://127.0.0.1:2379"
	path := "/config/test"

	viper.RemoteConfig = &utils.RemoteConfig{}

	ReConfig = viper.New()
	err := ReConfig.AddRemoteProvider("etcd", endpoint, path)
	if err!=nil {
		return err
	}
	ReConfig.SetConfigType("yaml")
	err = ReConfig.ReadRemoteConfig()
	if err!=nil {
		return err
	}
	err = ReConfig.WatchRemoteConfigOnChannel()
	if err!=nil {
		return err
	}
	return nil
}
