package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"
	"slg/GameServer"
	"slg/GateManager"
	"syscall"
)

var Type = flag.String("type", "auth", "Input Server Type")
var Group = flag.String("group", "king_war", "Input Server Type")
var AllowConnect = flag.Bool("allowConnect", false, "")

func main() {
	flag.Parse()
	f, _ := os.Create("./" + *Type + ".pprof")
	pprof.StartCPUProfile(f)
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("退出", s)
				pprof.StopCPUProfile()
			default:
				fmt.Println("other", s)
			}
		}
	}()

	switch *Type {
	case "auth":
		GateManager.Run(*Type, *Group, *AllowConnect)
		break
	case "game":
		GameServer.Run(*Type, *Group, *AllowConnect)
		break
	default:
		fmt.Println("please select service type")
		break
	}
}
