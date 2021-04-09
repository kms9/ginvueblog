package main

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
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

var ConfigMap = map[string]map[string]string{
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

func ConnNacos()  {
	tconfig:=ConfigMap["208"]
	var (
		serverAdd = tconfig["serverAdd"]
		dataId    = tconfig["dataId"]
		group     = tconfig["group"]
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
	//err = client.ListenConfig(vo.ConfigParam{
	//	DataId: dataId,
	//	Group:  group,
	//	OnChange: func(namespace, group, dataId, data string) {
	//		fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
	//	},
	//})
	//if err!=nil{
	//	log.Error(err)
	//}
}
func main()  {

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt,os.Kill)
	fmt.Println("start!")


	ConnNacos()

	s := <-c
	fmt.Println("stop,signal:",s)
}