package setup

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

var NacosConfig  *viper.Viper

// StartLogger 初始日志
func StartNConfig() error {

	var (
		serverAdd = "127.0.0.1"
		dataId    = "qq-config"
		group     = "test"
		//port      = 8848
	)

	//config := viper.New()
	sc := []constant.ServerConfig{
		{
			IpAddr: serverAdd,
			Port:   8848,
		},
	}
	//or a more graceful way to create ServerConfig
	_ = []constant.ServerConfig{
		*constant.NewServerConfig(
			serverAdd,
			8848,
			constant.WithScheme("http"),
			constant.WithContextPath("/nacos")),
	}

	cc := constant.ClientConfig{
		NamespaceId:         "e525eafa-f7d7-4029-83d9-008937f9d468", //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}
	//or a more graceful way to create ClientConfig
	_ = *constant.NewClientConfig(
		constant.WithNamespaceId("e525eafa-f7d7-4029-83d9-008937f9d468"),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithRotateTime("1h"),
		constant.WithMaxAge(3),
		constant.WithLogLevel("debug"),
	)

	// a more graceful way to create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}

	//publish config
	//config key=dataId+group+namespaceId
	//_, err = client.PublishConfig(vo.ConfigParam{
	//	DataId: dataId,
	//	Group:  group,
	//	Content: "hello world!",
	//})
	//_, err = client.PublishConfig(vo.ConfigParam{
	//	DataId: dataId,
	//	Group:  group,
	//	Content: "hello world!",
	//})
	//if err != nil {
	//	fmt.Printf("PublishConfig err:%+v \n", err)
	//}

	//get config
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err!=nil{
		Logger.Logger.Error(err)
	}
	fmt.Println("GetConfig,config :" + content)

	//Listen config change,key=dataId+group+namespaceId.
	err = client.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
		},
	})
	if err!=nil{
		Logger.Logger.Error(err)
	}
	err = client.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
		},
	})
	if err!=nil{
		Logger.Logger.Error(err)
	}
	//_, err = client.PublishConfig(vo.ConfigParam{
	//	DataId: dataId,
	//	Group:  group,
	//	Content: "test-listen",
	//})
	//
	//time.Sleep(2 * time.Second)
	//
	//_, err = client.PublishConfig(vo.ConfigParam{
	//	DataId: dataId,
	//	Group:  group,
	//	Content: "test-listen",
	//})
	//
	//time.Sleep(2 * time.Second)

	//cancel config change
	//err = client.CancelListenConfig(vo.ConfigParam{
	//	DataId: dataId,
	//	Group:  group,
	//})
	//
	//time.Sleep(2 * time.Second)
	//_, err = client.DeleteConfig(vo.ConfigParam{
	//	DataId: dataId,
	//	Group:  group,
	//})
	//time.Sleep(5 * time.Second)
	//
	//searchPage, _ := client.SearchConfig(vo.SearchConfigParm{
	//	Search:   "blur",
	//	DataId:   "",
	//	Group:    "",
	//	PageNo:   1,
	//	PageSize: 10,
	//})
	//fmt.Printf("Search config:%+v \n", searchPage)
	return nil
}