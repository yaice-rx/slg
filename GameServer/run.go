package GameServer

import (
	"SLGGAME/GameServer/Logic"
	"SLGGAME/GameServer/Package"
	"SLGGAME/Protocol/inside"
	"SLGGAME/Protocol/outside"
	"SLGGAME/Service"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/sirupsen/logrus"
	"github.com/yaice-rx/yaice"
	"github.com/yaice-rx/yaice/config"
	"github.com/yaice-rx/yaice/log"
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
	s.server.AddRouter(&outside.C2SGameCert{}, Logic.C2SGameCertHandler)
}

func (s *GameServer) BeforeRunThreadHook() {
	s.server.WatchServeNodeData(func(isAdd mvccpb.Event_EventType, config *config.Config) {
		switch isAdd {
		case mvccpb.PUT:
			fmt.Println("add", config)
			if config.TypeId == "auth" {
				conn := s.server.Dial(nil, "tcp", config.InHost+":"+strconv.Itoa(config.InPort))
				if conn == nil {
					log.AppLogger.Debug("connect error")
					return
				}
				conn.Send(&inside.RGameAuthRegisterRequest{Host: s.config.OutHost, Port: int32(s.config.OutPort)})
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
		data := s.server.Listen(Package.NewPackage(), "tcp", 30001, 30100)
		outerPort <- data
	}()
	s.config.OutPort = <-outerPort
	//开启内网
	s.config.InHost = "10.0.0.10"
	go func() {
		data := s.server.Listen(nil, "tcp", 30001, 30100)
		insidePort <- data
	}()
	s.config.InPort = <-insidePort
	//注册服务配置
	s.server.RegisterServeNodeData(*s.config)
	//关闭端口通道
	close(outerPort)
	close(insidePort)
	//开启连接
	for _, serverConf := range s.server.GetServeNodeData("") {
		if serverConf.TypeId == "auth" {
			conn := s.server.Dial(nil, "tcp", serverConf.InHost+":"+strconv.Itoa(serverConf.InPort))
			if conn == nil {
				log.AppLogger.Debug("connect error")
				return
			}
			conn.Send(&inside.RGameAuthRegisterRequest{Host: s.config.OutHost, Port: int32(s.config.OutPort)})
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
