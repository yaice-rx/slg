package AuthManager

import (
	"SLGGAME/AuthManager/Logic"
	"SLGGAME/AuthManager/ServiceGroup"
	"SLGGAME/Protocol/inside"
	"SLGGAME/Service"
	"github.com/yaice-rx/yaice"
	"github.com/yaice-rx/yaice/config"
	"github.com/yaice-rx/yaice/log"
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
	s.server.AddRouter(&inside.RGameAuthRegisterRequest{}, ServiceGroup.ServiceRegisterConn)

	s.server.AddRouter(&inside.RGameAuthPingRequest{}, ServiceGroup.ServicePingConn)
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
	s.config.OutPort = 50001
	go func() {
		http.HandleFunc("/login_dev", Logic.Login)
		http.ListenAndServe(":50001", nil)
	}()
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
}

func (s *AuthServer) ObserverPProf(addr string) {
	runtime.GOMAXPROCS(6)              // 限制 CPU 使用数，避免过载
	runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪
	runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪
	go func() {
		// 启动一个 http server，注意 pprof 相关的 handler 已经自动注册过了
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.AppLogger.Error(err.Error())
		}
		os.Exit(0)
	}()
}
