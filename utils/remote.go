package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/spf13/viper"
	"go.etcd.io/etcd/clientv3"
)

type RemoteConfig struct {
	viper.RemoteProvider

	Username string
	Password string
}

func (c *RemoteConfig) Get(rp viper.RemoteProvider) (io.Reader, error) {
	c.RemoteProvider = rp

	return c.get()
}

func (c *RemoteConfig) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	c.RemoteProvider = rp

	return c.get()
}

func (c *RemoteConfig) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	c.RemoteProvider = rp

	rr := make(chan *viper.RemoteResponse)
	stop := make(chan bool)

	go func() {
		for {
			client, err := c.newClient()

			if err != nil {
				time.Sleep(time.Duration(time.Second))
				fmt.Println("拉取刷新时间为", time.Second, "s" )
				continue
			}

			defer client.Close()

			ch := client.Watch(context.Background(), c.RemoteProvider.Path())

			select {
			case <-stop:
				return
			case res := <-ch:
				fmt.Println("触发更新逻辑 ",  )
				for _, event := range res.Events {
					fmt.Printf("触发更新逻辑, 更新为: %+v ", event)
					rr <- &viper.RemoteResponse{
						Value: event.Kv.Value,
					}
				}

			}
		}
	}()

	return rr, stop
}

func (c *RemoteConfig) newClient() (*clientv3.Client, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{c.Endpoint()},
		Username:  c.Username,
		Password:  c.Password,
	})

	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *RemoteConfig) get() (io.Reader, error) {
	client, err := c.newClient()

	if err != nil {
		return nil, err
	}

	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := client.Get(ctx, c.Path())
	cancel()

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(resp.Kvs[0].Value), nil
}