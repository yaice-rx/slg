package AuthManager

import (
	"SLGGAME/Protocol/inside"
	"SLGGAME/Service"
	"github.com/sirupsen/logrus"
	"github.com/yaice-rx/yaice"
	"github.com/yaice-rx/yaice/config"
	"github.com/yaice-rx/yaice/network"
	"net/http"
	"os"
	"runtime"
)

type AuthServer struct {
	type_     string
	groupName string
	config    *config.Config
	server    yaice.IServer
}

func NewServer(type_ string, serverGroup string) Service.IService {
	conf := new(config.Config)
	conf.TypeId = type_
	conf.ServerGroup = serverGroup
	s := &AuthServer{
		type_:     type_,
		groupName: serverGroup,
		config:    conf,
	}
	server := yaice.NewServer([]string{"127.0.0.1:2379"})
	s.server = server
	return s
}

func (s *AuthServer) RegisterProtoHandler() {
	//开启路由
	s.server.AddRouter(&inside.RGameAuthRegisterRequest{}, func(conn network.IConn, content []byte) {
		logrus.Info("register .....")
	})

	s.server.AddRouter(&inside.RGameAuthPingRequest{}, func(conn network.IConn, content []byte) {
		logrus.Info("ping .....")
	})
}

func (s *AuthServer) BeforeRunThreadHook() {

}

func (s *AuthServer) AfterRunThreadHook() {

}

func (s *AuthServer) Run() {
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
}

func (s *AuthServer) ObserverPProf(addr string) {
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
