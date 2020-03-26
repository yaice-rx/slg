package GameOuter

import (
	"SLGGAME/Protocol/outside"
	"SLGGAME/Token"
	"github.com/golang/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
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

func C2SGameCertHandler(conn network.IConn, content []byte) {
	log.AppLogger.Info("有服务器连接上来了:" + conn.GetConn().(*net.TCPConn).RemoteAddr().String())
	_ProtoData := outside.C2SGameCert{}
	err := proto.Unmarshal(content, &_ProtoData)
	if err != nil {
		return
	}
	tokenSignLen := _ProtoData.Token[0:2]
	tokenLen := utils.BytesToShort(tokenSignLen)
	var token Token.Token
	jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(_ProtoData.Token[2+tokenLen:], &token)
	data := []byte{}
	data = append(append(append(append(data, byte(2)), byte(1)), utils.LongToBytes(-24056111824897)...), utils.LongToBytes(_ProtoData.LoginSeq)...)
	msgData := append(append(utils.IntToBytes(8+1+1+8+8), utils.LongToBytes(int64(utils.GenerateCRCCheckCode(data)))...), data...)
	conn.SendByte(msgData)
}
