package setup

import (
	"errors"
	"sync"
	"github.com/kms9/publicyc/pkg/store/oredis"
)

var redisOnce sync.Once
var Cache *oredis.Redis

// startRedis 启动并初始化redis
func StartRedis() error {
	var err error
	redisOnce.Do(func() {
		Cache = oredis.UseConfig("ginblog").Build()
		if Cache == nil {
			err = errors.New("setup.redis start failed")
		}
	})
	return err
}
