package Logic

import (
	"SLGGAME/Protocol/outside"
	"github.com/golang/protobuf/proto"
	"github.com/yaice-rx/yaice/network"
	"github.com/yaice-rx/yaice/utils"
)

type cert struct {
	MsgType    int8
	AuthResult int8
	ServerGuid int64
	LoginSeq   int64
}

func C2SGameCertHandler(conn network.IConn, content []byte) {
	_ProtoData := outside.C2SGameCert{}
	err := proto.Unmarshal(content, &_ProtoData)
	if err != nil {
		return
	}
	data := []byte{}
	data = append(append(append(append(data, byte(2)), byte(1)), utils.LongToBytes(-24056111824897)...), utils.LongToBytes(_ProtoData.LoginSeq)...)
	msgData := append(append(utils.IntToBytes(8+1+1+8+8), utils.LongToBytes(int64(utils.GenerateCRCCheckCode(data)))...), data...)
	conn.SendByte(msgData)
}
