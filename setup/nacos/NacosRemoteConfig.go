package nacos

import (
	"bytes"
	"errors"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
	crypt "github.com/xordataexchange/crypt/config"
)


var (
	ErrUnsupportedProvider = errors.New("This configuration manager is not supported")

	_ viperConfigManager = nacosConfigManager{}
	// getConfigManager方法每次返回新对象导致缓存无效，
	// 这里通过endpoint作为key复一个对象
	// key: endpoint+appid value: agollo.Agollo
	agolloMap sync.Map
)


var (
	dataId string
	group string
	nameSpaceId string
	defaultConfigType = "yaml"
	defaultNacosOptions = vo.NacosClientParam{

	}

)

type nacosConfigManager struct {
	nConfigClient config_client.IConfigClient
}

func (n nacosConfigManager) Get( params vo.ConfigParam ) ([]byte, error) {
	content, err := n.nConfigClient.GetConfig(params)
	if err != nil {
		return nil, err
	}
	return []byte(content), nil
}

func (n nacosConfigManager) Watch(params vo.ConfigParam, stop chan bool) <-chan *RemoteResponse {
	resp := make(chan *viper.RemoteResponse)


	params.OnChange = func(namespace, group, dataId, data string) {

	}

	go func() {
		for {
			select {
			case <-stop:
				return
			case r := <-backendResp:
				if r.Error != nil {
					resp <- &viper.RemoteResponse{
						Value: nil,
						Error: r.Error,
					}
					continue
				}

				configType := getConfigType(namespace)
				value, err := marshalConfigs(configType, r.NewValue)

				resp <- &viper.RemoteResponse{Value: value, Error: err}
			}
		}
	}()
	return resp







	err :=  n.nConfigClient.ListenConfig(vo.ConfigParam{
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

}


type viperConfigManager interface {
	Get(key string) ([]byte, error)
	Watch(key string, stop chan bool) <-chan *viper.RemoteResponse
}

type configProvider struct {
}

func (rc configProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {
	cmt, err := getConfigManager(rp)
	if err != nil {
		return nil, err
	}

	var b []byte
	switch cm := cmt.(type) {
	case viperConfigManager:
		b, err = cm.Get(rp.Path())
	case crypt.ConfigManager:
		b, err = cm.Get(rp.Path())
	}

	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func (rc configProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	cmt, err := getConfigManager(rp)
	if err != nil {
		return nil, err
	}

	var resp []byte
	switch cm := cmt.(type) {
	case viperConfigManager:
		resp, err = cm.Get(rp.Path())
	case crypt.ConfigManager:
		resp, err = cm.Get(rp.Path())
	}

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(resp), nil
}

func (rc configProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	cmt, err := getConfigManager(rp)
	if err != nil {
		return nil, nil
	}

	switch cm := cmt.(type) {
	case viperConfigManager:
		quitwc := make(chan bool)
		viperResponsCh := cm.Watch(rp.Path(), quitwc)
		return viperResponsCh, quitwc
	default:
		ccm := cm.(crypt.ConfigManager)
		quit := make(chan bool)
		quitwc := make(chan bool)
		viperResponsCh := make(chan *viper.RemoteResponse)
		cryptoResponseCh := ccm.Watch(rp.Path(), quit)
		// need this function to convert the Channel response form crypt.Response to viper.Response
		go func(cr <-chan *crypt.Response, vr chan<- *viper.RemoteResponse, quitwc <-chan bool, quit chan<- bool) {
			for {
				select {
				case <-quitwc:
					quit <- true
					return
				case resp := <-cr:
					vr <- &viper.RemoteResponse{
						Error: resp.Error,
						Value: resp.Value,
					}

				}

			}
		}(cryptoResponseCh, viperResponsCh, quitwc, quit)

		return viperResponsCh, quitwc
	}
}

func newNacosConfigManager(nacosClientParam  vo.NacosClientParam)(*nacosConfigManager, error)  {
	if dataId == "" {
		return nil, errors.New("The dataId is not set")
	}

	if group == "" {
		return nil, errors.New("The group is not set")
	}

	if nameSpaceId == "" {
		return nil, errors.New("The nameSpaceId is not set")
	}

	client, err := clients.NewConfigClient(
		nacosClientParam,
	)

	if err != nil {
		return nil, err
	}

	return &nacosConfigManager{
		nConfigClient: client,
	}, nil
}

func getConfigManager(rp viper.RemoteProvider) (interface{}, error) {
	if rp.SecretKeyring() != "" {
		kr, err := os.Open(rp.SecretKeyring())
		if err != nil {
			return nil, err
		}
		defer kr.Close()

		switch rp.Provider() {
		case "etcd":
			return crypt.NewEtcdConfigManager([]string{rp.Endpoint()}, kr)
		case "consul":
			return crypt.NewConsulConfigManager([]string{rp.Endpoint()}, kr)
		case "nacos":
			return nil, errors.New("The Nacos configuration manager is not support encrypted")
		default:
			return nil, ErrUnsupportedProvider
		}
	} else {
		switch rp.Provider() {
		case "etcd":
			return crypt.NewStandardEtcdConfigManager([]string{rp.Endpoint()})
		case "consul":
			return crypt.NewStandardConsulConfigManager([]string{rp.Endpoint()})
		case "nacos":
			return newNacosConfigManager(defaultNacosOptions)
		default:
			return nil, ErrUnsupportedProvider
		}
	}
}

func init() {
	viper.SupportedRemoteProviders = append(
		viper.SupportedRemoteProviders,
		"nacos",
	)
	viper.RemoteConfig = &configProvider{}
}
