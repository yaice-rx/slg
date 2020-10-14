package main

import (
	"SLGGAME/AuthManager"
	"SLGGAME/GameServer"
	"SLGGAME/RouterManager"
	"SLGGAME/Service"
	"flag"
	"fmt"
	_ "net/http/pprof"
)

var Type 	= flag.String("type", "auth", "Input Server Type")
var Group 	= flag.String("group", "king_war", "Input Server Type")

func main() {
	flag.Parse()
	IsEndRuning := make(chan bool)
	defer func() {
		//关闭服务通道
		close(IsEndRuning)
	}()
	var serve Service.IService
	switch *Type {
	case "auth":
		serve = AuthManager.NewServer(*Type, *Group)
		serve.ObserverPProf(":1592")
		break
	case "game":
		serve = GameServer.NewServer(*Type, *Group)
		serve.ObserverPProf(":1593")
		break
	case "router":
		serve = RouterManager.NewServer(*Type, *Group)
		serve.ObserverPProf(":1594")
		break
	default:
		fmt.Println("please select service type")
		break
	}
	if serve == nil {
		IsEndRuning <- true
	}
	//注册协议
	serve.RegisterProtoHandler()
	//启动服务之前的操作
	serve.BeforeRunThreadHook()
	//启动服务
	serve.Run()
	//启动之后的操作
	serve.AfterRunThreadHook()
	//结束进程
	<-IsEndRuning
	return
}
