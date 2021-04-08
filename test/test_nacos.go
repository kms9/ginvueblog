package main

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"
	"os"
	"fmt"
)

func init()  {
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})
		// Output to stdout instead of the default stderr
		// Can be any io.Writer, see below for File example
		log.SetOutput(os.Stdout)

		// Only log the warning severity or above.
		log.SetLevel(log.TraceLevel)
		log.SetReportCaller(true)
}
func ConnNacos()  {
	var (
		serverAdd = "127.0.0.1"
		dataId    = "qq-config"
		group     = "DEFAULT_GROUP"
		nameSpaceId = "5a091299-14bd-459b-95ee-78a8d770e65e"
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
	}


	// a more graceful way to create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		log.Error(err)
		panic(err)
	}

	//get config
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err!=nil{
		log.Error(err)
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
		log.Error(err)
	}
	err = client.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
		},
	})
	if err!=nil{
		log.Error(err)
	}
}
func main()  {
	ConnNacos()
}