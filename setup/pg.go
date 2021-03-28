package setup

import (
	"errors"
	"sync"
	"github.com/kms9/publicyc/pkg/store/opg"
)

var pgOnce sync.Once
var DB *opg.DB

// StartDB 启动并初始化数据库
func StartDB() error {
	var err error
	pgOnce.Do(func() {
		DB = opg.UseConfig("ginblog").Build()
		if DB == nil {
			err = errors.New("setup.db start failed")
		}
	})

	return err
}
