package GameServer

import (
	"SLGGAME/Protocol/inside"
	"SLGGAME/Service"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/sirupsen/logrus"
	"github.com/yaice-rx/yaice"
	"github.com/yaice-rx/yaice/config"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

type GameServer struct {
	type_     string
	groupName string
	config    *config.Config
	server    yaice.IServer
}

func NewServer(type_ string, serverGroup string) Service.IService {
	conf := new(config.Config)
	conf.TypeId = type_
	conf.ServerGroup = serverGroup
	s := &GameServer{
		type_:     type_,
		groupName: serverGroup,
		config:    conf,
	}
	server := yaice.NewServer([]string{"127.0.0.1:2379"})
	s.server = server
	return s
}

func (s *GameServer) RegisterProtoHandler() {

}

func (s *GameServer) BeforeRunThreadHook() {
	s.server.WatchNodeData(func(isAdd mvccpb.Event_EventType, config *config.Config) {
		switch isAdd {
		case mvccpb.PUT:
			fmt.Println("add")
			if config.TypeId == "auth" {
				conn := s.server.DialTCP(config.InHost + ":" + strconv.Itoa(config.InPort))
				conn.Send(&inside.RGameAuthRegisterRequest{Pid: strconv.Itoa(os.Getpid())})
				go func() {
					for _ = range time.Tick(5 * time.Second) {
						conn.Send(&inside.RGameAuthPingRequest{})
					}
				}()
			}
			break
		case mvccpb.DELETE:
			fmt.Println("delete", config)
			break
		}
	})
}

func (s *GameServer) AfterRunThreadHook() {

}

func (s *GameServer) Run() {
	//设置监听端口通道
	outerPort := make(chan int)
	insidePort := make(chan int)
	//开启外网
	s.config.OutHost = "10.0.0.10"
	go func() {
		data := s.server.ListenTCP(30001, 30100)
		outerPort <- data
	}()
	s.config.OutPort = <-outerPort
	//开启内网
	s.config.InHost = "10.0.0.10"
	go func() {
		data := s.server.ListenTCP(30001, 30100)
		insidePort <- data
	}()
	s.config.InPort = <-insidePort
	//关闭端口通道
	close(outerPort)
	close(insidePort)
	//注册服务配置
	s.server.RegisterNodeData(*s.config)
	//开启连接
	for _, serverConf := range s.server.GetNodeData("") {
		if serverConf.TypeId == "auth" {
			conn := s.server.DialTCP(serverConf.InHost + ":" + strconv.Itoa(serverConf.InPort))
			conn.Send(&inside.RGameAuthRegisterRequest{Pid: strconv.Itoa(os.Getpid())})
			go func() {
				for _ = range time.Tick(5 * time.Second) {
					conn.Send(&inside.RGameAuthPingRequest{})
				}
			}()
		}
	}
}

func (s *GameServer) ObserverPProf(addr string) {
	runtime.GOMAXPROCS(6)              // 限制 CPU 使用数，避免过载
	runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪
	runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪
	go func() {
		// 启动一个 http server，注意 pprof 相关的 handler 已经自动注册过了
		if err := http.ListenAndServe(addr, nil); err != nil {
			logrus.Debug(err)
		}
		os.Exit(0)
	}()
}
