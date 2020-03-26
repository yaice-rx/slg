package ServiceGroup

import (
	"SLGGAME/Protocol/inside"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/yaice-rx/yaice/log"
	"github.com/yaice-rx/yaice/network"
)

func RegisterConnHandler(conn network.IConn, content []byte) {
	data := inside.RGameAuthRegisterRequest{}
	err := proto.Unmarshal(content, &data)
	if err != nil {
		log.AppLogger.Debug("连接服务器的参数错误，不能解析。。。" + err.Error())
		return
	}
}

func PingConnHandler(conn network.IConn, content []byte) {
	logrus.Info("ping .....")
}

func LoginHandler(conn network.IConn, content []byte) {
	data := inside.RGameAuthLoginRequest{}
	err := proto.Unmarshal(content, &data)
	if err != nil {
		log.AppLogger.Debug("玩家登陆验证数据错误，不能解析。。。" + err.Error())
		return
	}
	conn.Send(&inside.RGameAuthLoginCallback{State: true})
}
