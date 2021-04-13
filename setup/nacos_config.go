package setup

import (

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"strings"
	"fmt"
)

var NacosConfig  *viper.Viper

// StartLogger 初始日志
func StartNacosConfig() error {
	err := ConnNacos()
	if err!=nil {
		Logger.Logger.Info(err)
		return err
	}
	return nil
}

var ConfigMap = map[string]map[string]string{
	"208" : {
		"serverAdd" : "10.8.8.208",
		"dataId"    : "qq-redpackets",
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

func ConnNacos()  error {
	tconfig:=ConfigMap["201"]
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

	// a more graceful way to create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		//log.Error(err)
		panic(err)
		return err
	}

	//get config
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err!=nil{
		//log.Error(err)
		return err
	}
	fmt.Println("GetConfig,config :" + content)
	NacosConfig = viper.New()

	Logger.Logger.Info(NacosConfig)

	NacosConfig.SetConfigType("yaml")
	err = NacosConfig.ReadConfig(strings.NewReader(content))
	if err != nil {
		fmt.Println("Viper解析配置失败:", err)
		return err
	}

	Logger.Logger.Info(NacosConfig)

	//Listen config change,key=dataId+group+namespaceId.
	err = client.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
			err = NacosConfig.ReadConfig(strings.NewReader(data))
			Logger.Logger.Info(NacosConfig)
			if err != nil {
				fmt.Println("Viper解析配置失败:", err)
				panic(err)
				return
			}
		},
	})
	if err!=nil{
		//log.Error(err)
		Logger.Logger.Info(err)
		return err
	}
	return nil
}