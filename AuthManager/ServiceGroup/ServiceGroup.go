package ServiceGroup

import (
	"SLGGAME/Protocol/inside"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/yaice-rx/yaice/log"
	"github.com/yaice-rx/yaice/network"
)

type Service struct {
	Host  string
	Port  int32
	group int
}

type Services struct {
	group map[int]Service
}

var ServicesMgr = constructServices()

func constructServices() Services {
	return Services{
		group: make(map[int]Service),
	}
}

func (s *Services) AddService(data Service) {
	s.group[data.group] = data
}

func (s *Services) GetService(gid int) Service {
	return s.group[gid]
}

func ServiceRegisterConn(conn network.IConn, content []byte) {
	data := inside.RGameAuthRegisterRequest{}
	err := proto.Unmarshal(content, &data)
	if err != nil {
		log.AppLogger.Debug("连接服务器的参数错误，不能解析。。。" + err.Error())
		return
	}
	log.AppLogger.Info("服务器：" + data.Host + ",已经连接上来")
	ServicesMgr.AddService(Service{Host: data.Host, Port: data.Port, group: 1})
}

func ServicePingConn(conn network.IConn, content []byte) {
	logrus.Info("ping .....")
}
