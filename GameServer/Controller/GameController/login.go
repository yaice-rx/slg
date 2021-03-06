package GameController

import (
	"SLGGAME/GameServer/Session"
	"SLGGAME/Protocol/inside"
	"SLGGAME/Protocol/outside"
	"SLGGAME/Service"
	"github.com/golang/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/yaice-rx/yaice/log"
	"github.com/yaice-rx/yaice/network"
	"github.com/yaice-rx/yaice/utils"
	"net"
)

type cert struct {
	MsgType    int8
	AuthResult int8
	ServerGuid int64
	LoginSeq   int64
}

var TokenLoginData []byte

func PlayerLoginHandler(conn network.IConn, content []byte) {
	//添加在线玩家列表
	log.AppLogger.Info("有服务器连接上来了:" + conn.GetConn().(*net.TCPConn).RemoteAddr().String())
	_ProtoData := outside.C2SGameCert{}
	err := proto.Unmarshal(content, &_ProtoData)
	if err != nil {
		return
	}
	tokenSignLen := _ProtoData.Token[0:2]
	tokenLen := utils.BytesToShort(tokenSignLen)
	var token Service.Token
	jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(_ProtoData.Token[2+tokenLen:], &token)
	//将玩家加入在线列表
	Session.PlayerContainsGameMgr.Add(token.Guid, conn)
	data := []byte{}
	data = append(append(append(append(data, byte(2)), byte(1)), utils.LongToBytes(-24056111824897)...), utils.LongToBytes(_ProtoData.LoginSeq)...)
	msgData := append(append(utils.IntToBytes(8+1+1+8+8), utils.LongToBytes(int64(utils.GenerateCRCCheckCode(data)))...), data...)
	//todo 向Auth服务器请求连接，通过玩家的请求的消息中获取对应的auth服务器的guid
	severConn := Session.AuthContainsGameMgr.Get(token.SessionId)
	severConn.Send(&inside.RGameAuthLoginRequest{PlayerGuid: token.Guid})
	TokenLoginData = msgData
}

func PlayerRegisterHandler(conn network.IConn, content []byte) {
	logrus.Info("玩家：", conn.GetGuid(), " ,已经注册成功")
}
