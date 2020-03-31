package GameServer

import (
	"SLGGAME/GameServer/Controller/AuthController"
	"SLGGAME/GameServer/Controller/GameController"
	"SLGGAME/GameServer/Package"
	"SLGGAME/GameServer/Session"
	"SLGGAME/Protocol/inside"
	"SLGGAME/Protocol/outside"
	"SLGGAME/Service"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/sirupsen/logrus"
	"github.com/yaice-rx/yaice"
	"github.com/yaice-rx/yaice/config"
	"github.com/yaice-rx/yaice/log"
	"github.com/yaice-rx/yaice/network"
	"github.com/yaice-rx/yaice/utils"
	"net/http"
	"os"
	"runtime"
	"strconv"
)

type GameServer struct {
	type_     string
	groupName string
	config    *config.Config
	server    yaice.IServer
}

var ServerConfigMgr *config.Config

func NewServer(type_ string, serverGroup string) Service.IService {
	conf := new(config.Config)
	conf.TypeId = type_
	conf.ServerGroup = serverGroup
	conf.Pid = utils.GenSonyflake()
	s := &GameServer{
		type_:     type_,
		groupName: serverGroup,
		config:    conf,
	}
	server := yaice.NewServer([]string{"127.0.0.1:2379"})
	s.server = server
	ServerConfigMgr = conf
	return s
}

func (s *GameServer) RegisterProtoHandler() {
	//auth注册回调
	s.server.AddRouter(&inside.RGameAuthRegisterCallback{}, AuthController.AuthTGameRegisterResultFunc)
	//auth服ping回调
	s.server.AddRouter(&inside.RGameAuthPingCallback{}, AuthController.AuthTGamePingResultFunc)
	//auth服player回调
	s.server.AddRouter(&inside.RGameAuthLoginCallback{}, AuthController.AuthTGameLoginResultFunc)
	//玩家登陆
	s.server.AddRouter(&outside.C2SGameCert{}, GameController.PlayerLoginHandler)
	//网络注册
	s.server.AddRouter(&outside.C2GRegister{}, GameController.PlayerRegisterHandler)
}

func (s *GameServer) BeforeRunThreadHook() {
	s.server.WatchServeNodeData(func(isAdd mvccpb.Event_EventType, key []byte, value config.IConfig) {
		switch isAdd {
		case mvccpb.PUT:
			if value.GetTypeId() == "auth" && Session.AuthContainsGameMgr.Get(value.GetPid()) == nil {
				conn := s.server.Dial(nil, "tcp", value.GetInHost()+":"+strconv.Itoa(value.GetInPort()),
					network.WithMax(10))
				if conn == nil {
					log.AppLogger.Debug("connect error")
					return
				}
				conn.Send(&inside.RGameAuthRegisterRequest{
					Host: s.config.OutHost,
					Port: int32(s.config.OutPort),
					Guid: config.ConfInstance().GetPid(),
				})
			}
			break
		case mvccpb.DELETE:
			keyStr := string(key)
			if keyStr != "" {
				//keyMap := strings.Split(keyStr,"\\")
				//Ikey,_ := strconv.ParseUint(keyMap[3], 10, 64)
			}
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
	defer func() {
		//关闭端口通道
		close(outerPort)
		close(insidePort)
	}()
	//开启外网
	ServerConfigMgr.OutHost = "127.0.0.1"
	go func() {
		data := s.server.Listen(Package.NewPackage(), "tcp", 30001, 30100,
			func(conn interface{}) bool {
				return true
			})
		outerPort <- data
	}()
	ServerConfigMgr.OutPort = <-outerPort
	//开启内网
	ServerConfigMgr.InHost = "127.0.0.1"
	go func() {
		data := s.server.Listen(nil, "tcp", 30001, 30100,
			func(conn interface{}) bool {
				return true
			})
		insidePort <- data
	}()
	ServerConfigMgr.InPort = <-insidePort
	//注册服务配置
	s.server.RegisterServeNodeData()
	//开启连接
	for _, value := range s.server.GetServeNodeData("") {
		if value.GetTypeId() == "auth" && Session.AuthContainsGameMgr.Get(value.GetPid()) == nil {
			conn := s.server.Dial(nil, "tcp", value.GetInHost()+":"+strconv.Itoa(value.GetInPort()),
				network.WithMax(10))
			if conn == nil {
				log.AppLogger.Debug("connect error")
				return
			}
			//注册信息
			conn.Send(&inside.RGameAuthRegisterRequest{
				Host: s.config.OutHost,
				Port: int32(s.config.OutPort),
				Guid: config.ConfInstance().GetPid(),
			})
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
