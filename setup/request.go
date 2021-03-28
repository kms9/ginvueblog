package setup

import (
	"errors"
	"sync"

	"github.com/kms9/publicyc/pkg/client/ohttp"
)

var requestOnce sync.Once
var Request *ohttp.Requests

// StartRequest 启动并初始化请求
func StartRequest() error {
	var err error
	requestOnce.Do(func() {
		Request = ohttp.UseConfig("request").Build()
		if Request == nil {
			err = errors.New("setup.request start failed")
		}
	})
	return err
}
