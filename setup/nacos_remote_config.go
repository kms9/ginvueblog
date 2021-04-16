package setup

import (
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	nacosRemote "ginvueblog/setup/nacos"
)

var NacosRConfig  *viper.Viper

var DefaultConfigMap = map[string]map[string]string{
	"208" : {
		"serverAdd" : "10.8.8.208",
		"dataId"    : "qq-config",
		"group"     : "DEFAULT_GROUP",
		"nameSpaceId" : "2baea186-51ba-459e-bfa8-4c222d16c308",
	},

	"201": {
		"serverAdd" : "192.168.31.201",
		"dataId"    : "qq-config",
		"group"     : "DEFAULT_GROUP",
		"nameSpaceId" : "4793f393-e460-43df-bace-99c8ba4cbe06",
	},
}

// StartLogger 初始日志
func StartNacosRConfig() error {
	endpoint := "http://127.0.0.1:2379"
	path := "/config/test"
	tconfig:=DefaultConfigMap["208"]
	NacosRConfig  = viper.New()

	var (
		serverAdd 	= tconfig["serverAdd"]
		dataId    	= tconfig["dataId"]
		group     	= tconfig["group"]
		nameSpaceId = tconfig["nameSpaceId"]
		//port      = 8848
	)

	//config := viper.New()
	sc := []constant.ServerConfig{
		{
			IpAddr: serverAdd,
			Port:   8848,
			Scheme: "http",
			ContextPath: "/nacos",
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         nameSpaceId, //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "./config/log",
		CacheDir:            "./config/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
		//ListenInterval:
	}

	params:=vo.NacosClientParam{
		ClientConfig:  &cc,
		ServerConfigs: sc,
	}

	nacosRemote.SetDataID(dataId)
	nacosRemote.SetGroup(group)

	nacosRemote.SetNacosOptions(params)
	NacosRConfig.SetConfigType("yaml")

	//Logger.Logger.Info("nacosRemote set config")

	err:= NacosRConfig.AddRemoteProvider("nacos", endpoint, path)
	if err!=nil {
		Logger.Logger.Info("NacosRConfig.AddRemoteProvider err: ", err)
		return err
	}
	Logger.Logger.Info(viper.SupportedRemoteProviders)
	err = NacosRConfig.ReadRemoteConfig()
	if err!=nil {
		Logger.Logger.Info("NacosRConfig.ReadRemoteConfig err: ", err)
		return err
	}

	err = NacosRConfig.WatchRemoteConfigOnChannel()
	if err!=nil {
		return err
	}

	return nil
}



