package setup

import (
	"errors"
	"sync"
	"github.com/kms9/publicyc/pkg/store/omysql"
)

var mysqlOnce sync.Once
var MysqlDB *omysql.DB

// StartDB 启动并初始化数据库
func StartMysqlDB() error {
	var err error
	mysqlOnce.Do(func() {
		MysqlDB = omysql.UseConfig("ginblog").Build()
		if MysqlDB == nil {
			err = errors.New("setup.db start failed")
		}
	})

	return err
}
