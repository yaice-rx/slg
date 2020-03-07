package Package

import (
	"SLGGAME/Protocol/outside"
	"fmt"
	"github.com/yaice-rx/yaice/network"
	"github.com/yaice-rx/yaice/utils"
)

type Package struct {
}

const (
	MsgLengthLen    = 4
	MsgSumLen       = 8
	MsgTypeLen      = 1
	MsgPosLen       = 8
	MsgLoginSeqLen  = 8
	MsgTokenSignLen = 4
)

func NewPackage() network.IPacket {
	return &Package{}
}

func (p *Package) Pack(msg network.IMessage) []byte {
	fmt.Println("--------------------发送网络数据调用----------------")
	return nil
}

func (p *Package) Unpack(buff []byte) ([]byte, []byte, int32, error) {
	msgType := buff[MsgLengthLen+MsgSumLen : MsgLengthLen+MsgSumLen+MsgTypeLen]
	switch int(msgType[0]) {
	case 1:
		//当消息类型处于1的时候，利用protobuf伪造一个消息
		msgName := utils.GetProtoName(&outside.C2SGameCert{})
		protocolNum := utils.ProtocalNumber(msgName)
		return []byte{}, buff[MsgLengthLen+MsgSumLen+MsgTypeLen+MsgPosLen : MsgLengthLen+MsgSumLen+MsgTypeLen+MsgPosLen+MsgLoginSeqLen], protocolNum, nil
		break
	case 3:

		break
	}
	return nil, nil, 0, nil
}
