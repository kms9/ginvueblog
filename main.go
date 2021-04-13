package main

import (
	"fmt"
	"ginvueblog/cache"
	"ginvueblog/routes"
	"ginvueblog/services"
	"ginvueblog/setup"
	"github.com/kms9/publicyc"
	"github.com/kms9/publicyc/pkg/server/ogin"
	//_ "github.com/spf13/viper/remote"
)


// Engine ..
type Engine struct {
	yc.Application
}

// NewEngine 初始化相关方法
func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Start(
		setup.StartLogger,
		setup.StartRedis,
		//setup.StartConfig,
		setup.StartNacosConfig,
		//setup.StartDB,
		//setup.StartMysqlDB,
		//setup.StartRequest,
		cache.Init,
		services.ServiceInit,
		eng.serveHTTP,
	); err != nil {
		//onion_log.Panicf("Engine: %s", err)
		err := fmt.Errorf("Engine: %s", err)
		panic(err)
	}
	return eng
}

// serveHTTP 启动http
func (eng *Engine) serveHTTP() error {
	server := ogin.UseConfig("http").Build()
	routes.StartHttp(server)
	return eng.Serve(server)
}

func main() {
	//====old code
	// 引用数据库
	//upload.InitDb()
    // 引入路由组件
    //routes.InitRouter()
	//================old code end

	//config.Init()
	//model.Init()


	eng := NewEngine()
	if err := eng.Run(); err != nil {
		fmt.Println(err.Error())
	}
}
