package Logic

import (
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
	data := []byte{}
	data = append(append(append(append(data, byte(2)), byte(1)), utils.LongToBytes(-24056111824897)...), content...)
	msgData := []byte{}
	msgData = append(append(append(msgData, utils.IntToBytes(8+1+1+8+8)...), utils.LongToBytes(int64(utils.GenerateCRCCheckCode(data)))...), data...)
	conn.SendByte(msgData)
}
