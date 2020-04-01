package AuthController

import (
	"SLGGAME/GameServer/Controller/GameController"
	"SLGGAME/GameServer/Session"
	"SLGGAME/Protocol/inside"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/yaice-rx/yaice/network"
	"time"
)

func AuthTGameRegisterResultFunc(conn network.IConn, content []byte) {
	data := &inside.RGameAuthRegisterCallback{}
	proto.Unmarshal(content, data)
	go func() {
		for _ = range time.Tick(5 * time.Second) {
			conn.Send(&inside.RGameAuthPingRequest{})
		}
	}()
}

func AuthTGamePingResultFunc(conn network.IConn, content []byte) {
	logrus.Info("=========logic auth服务ping通回调")
}

func AuthTGameLoginResultFunc(conn network.IConn, content []byte) {
	logrus.Info("玩家登陆验证成功")
	data := &inside.RGameAuthLoginCallback{}
	proto.Unmarshal(content, data)
	playerConn := Session.PlayerContainsGameMgr.Get(data.Guid)
	playerConn.SendByte(GameController.TokenLoginData)
}
