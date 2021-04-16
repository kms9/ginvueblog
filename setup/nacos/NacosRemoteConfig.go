package nacos

import (
	"bytes"
	"errors"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/silenceper/log"
	"github.com/spf13/viper"
	crypt "github.com/xordataexchange/crypt/config"
	"io"
	"os"
)

var (
	ErrUnsupportedProvider = errors.New("This configuration manager is not supported")
	_ viperConfigManager = nacosConfigManager{}
)

func SetDataID(dataIdStr string) {
	dataId = dataIdStr
}

func SetGroup(groupName string) {
	group = groupName
}

func SetNacosOptions(params vo.NacosClientParam) {
	defaultNacosOptions = params
}

var (
	dataId string
	group  string
	defaultNacosOptions = vo.NacosClientParam{}
)

type nacosConfigManager struct {
	nConfigClient config_client.IConfigClient
}

func (n nacosConfigManager) Get(path string) ([]byte, error) {

	content, err := n.nConfigClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return []byte(content), nil
}

func (n nacosConfigManager) Watch(appId string, stop chan bool) <-chan *viper.RemoteResponse {
	resp := make(chan *viper.RemoteResponse)
	params := vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	}
	params.OnChange = func(namespace, group, dataId, data string) {
		resp <- &viper.RemoteResponse{
			Value: []byte(data),
			Error: nil,
		}
	}

	err := n.nConfigClient.ListenConfig(params)
	if err != nil {
		log.Error(err.Error())
		panic(err)
		return nil
	}

	go func(client config_client.IConfigClient, dataId, group string) {
		for {
			select {
			case <-stop:
				err := n.nConfigClient.CancelListenConfig(vo.ConfigParam{
					DataId: dataId,
					Group:  group,
				})
				if err != nil {
					log.Error(err.Error())
					panic(err)
					return
				}
				return
			}
		}
	}(n.nConfigClient, dataId, group)

	return resp
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
		log.Error(err.Error())
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
		log.Error(err.Error())
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func (rc configProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	cmt, err := getConfigManager(rp)
	if err != nil {
		log.Error(err.Error())
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
		log.Error(err.Error())
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

//nacosClientParam  vo.NacosClientParam
func newNacosConfigManager(endpoint string) (*nacosConfigManager, error) {
	if dataId == "" {
		return nil, errors.New("The dataId is not set")
	}

	if group == "" {
		return nil, errors.New("The group is not set")
	}

	client, err := newNacos()

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &nacosConfigManager{
		nConfigClient: client,
	}, nil
}

func newNacos() (config_client.IConfigClient, error) {
	nacosClient, err := clients.NewConfigClient(
		defaultNacosOptions,
	)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return nacosClient, nil
}

func getConfigManager(rp viper.RemoteProvider) (interface{}, error) {
	if rp.SecretKeyring() != "" {
		kr, err := os.Open(rp.SecretKeyring())
		if err != nil {
			return nil, err
		}
		defer func() {
			err = kr.Close()
			panic(err)
		}()

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
			return newNacosConfigManager(rp.Endpoint())
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
