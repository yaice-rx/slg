package AuthManager

import (
	"SLGGAME/AuthManager/Logic"
	"SLGGAME/AuthManager/ServiceGroup"
	"SLGGAME/Protocol/inside"
	"SLGGAME/Service"
	"SLGGAME/Token"
	jsoniter "github.com/json-iterator/go"
	"github.com/yaice-rx/yaice"
	"github.com/yaice-rx/yaice/config"
	"github.com/yaice-rx/yaice/log"
	"github.com/yaice-rx/yaice/utils"
	"net/http"
	"os"
	"runtime"
	"time"
)

type AuthServer struct {
	type_     string
	groupName string
	confMgr   config.IConfig
	server    yaice.IServer
}

func NewServer(type_ string, serverGroup string) Service.IService {
	s := &AuthServer{
		type_:     type_,
		groupName: serverGroup,
		confMgr:   config.ConfInstance(),
	}
	s.confMgr.SetTypeId(type_)
	s.confMgr.SetServerGroup(serverGroup)
	s.confMgr.SetPid(utils.GenSonyflake())
	server := yaice.NewServer([]string{"127.0.0.1:2379"})
	s.server = server
	return s
}

func (s *AuthServer) RegisterProtoHandler() {
	//服务内部注册
	s.server.AddRouter(&inside.RGameAuthRegisterRequest{}, ServiceGroup.RegisterConnHandler)
	//Ping
	s.server.AddRouter(&inside.RGameAuthPingRequest{}, ServiceGroup.PingConnHandler)
	//玩家登陆验证
	s.server.AddRouter(&inside.RGameAuthLoginRequest{}, ServiceGroup.LoginHandler)
}

func (s *AuthServer) BeforeRunThreadHook() {

}

func (s *AuthServer) AfterRunThreadHook() {

}

func (s *AuthServer) Run() {
	//设置监听端口通道
	outerPort := make(chan int)
	insidePort := make(chan int)
	//关闭端口通道
	defer func() {
		close(outerPort)
		close(insidePort)
	}()
	//开启外网
	s.confMgr.SetOutHost("10.0.0.10")
	s.confMgr.SetOutPort(50001)
	go func() {
		http.HandleFunc("/login_dev", s.Login)
		http.ListenAndServe(":50001", nil)
	}()
	//开启内网
	s.confMgr.SetInHost("10.0.0.10")
	go func() {
		data := s.server.Listen(nil, "tcp", 30001, 30100)
		insidePort <- data
	}()
	s.confMgr.SetInPort(<-insidePort)
	//注册服务配置
	s.server.RegisterServeNodeData()
	return
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

func (s *AuthServer) Login(w http.ResponseWriter, r *http.Request) {
	result := Token.Token{
		Guid:      time.Now().Unix(),
		SessionId: utils.GenSonyflake(),
		Port:      s.confMgr.GetInPort(),
		Host:      s.confMgr.GetInHost(),
		Result:    1,
	}
	data, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(result)
	w.Write(Logic.BuildToken(data))
}
