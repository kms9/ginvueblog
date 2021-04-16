package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

func SendMsgByTimer(test func(string))  {
	ticker:=time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Ticker tick.")
		test(time.Now().String())
	}
}

func GetMsg()  {
	//go func() {
	//	for {
	//		select {
	//		case r := <-backendResp:
	//			if r.Error != nil {
	//				resp <- &viper.RemoteResponse{
	//					Value: nil,
	//					Error: r.Error,
	//				}
	//				continue
	//			}
	//
	//			configType := getConfigType(namespace)
	//			value, err := marshalConfigs(configType, r.NewValue)
	//
	//			resp <- &viper.RemoteResponse{Value: value, Error: err}
	//		}
	//	}
	//}()

	for  {
		for v:=range eventBus{
			fmt.Println(v)
		}
	}

}
var eventBus =  make(chan string)

func main()  {
	//测试不使用定时器的消息主动推送逻辑

	//定时发送消息
	testFunc:= func(eb string) {
		//fmt.Println(time.Now())
		eventBus <- eb
	}

	go SendMsgByTimer(testFunc)
	go GetMsg()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt,os.Kill)
	s := <-c
	fmt.Println("stop,signal:",s)

}