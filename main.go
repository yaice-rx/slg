package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"slg/GameServer"
	"slg/GateManager"
)

var Type = flag.String("type", "auth", "Input Server Type")
var Group = flag.String("group", "king_war", "Input Server Type")
var AllowConnect = flag.Bool("allowConnect", false, "")

func main() {
	flag.Parse()
	/*runtime.GOMAXPROCS(6) // 限制 CPU 使用数，避免过载
	runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪
	runtime.SetBlockProfileRate(1) // 开启对阻塞操作的跟踪
	go func() {
		// 启动一个 http server，注意 pprof 相关的 handler 已经自动注册过了
		if err := http.ListenAndServe(":6060", nil); err != nil {
			logrus.Debug(err)
		}
		os.Exit(0)
	}()*/
	/*f, _ := os.Create("./" + *Type + ".pprof")
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
	}()*/

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
