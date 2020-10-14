package Package

import (
	"SLGGAME/CoreManager/NPackage"
	"SLGGAME/Protocol/outside"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/yaice-rx/yaice/network"
	"github.com/yaice-rx/yaice/utils"
)

type Package struct {
}

const (
	MsgLengthLen = 4
	MsgSumLen    = 8
	MsgTypeLen   = 1
	MsgIsPos     = 8
)

func NewPackage() network.IPacket {
	return &Package{}
}

func (p *Package) GetHeadLen() uint32 {
	return MsgLengthLen
}

func (p *Package) Pack(msg network.TransitData) []byte {
	return nil
}

func (p *Package) Unpack(binaryData []byte) (network.IMessage, error, func(conn network.IConn)) {
	//获取MsgType
	type_ := binaryData[MsgSumLen : MsgSumLen+MsgTypeLen][0]
	verifyData := binaryData[MsgSumLen:]
	verifySumCode := NPackage.GenerateCRCCheckCode(verifyData)
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)
	switch type_ {
	case 1:
		msg := &AuthMessage{}
		msg.MsgType = type_
		//读取sum
		msg.Sum = utils.BytesToLong(binaryData[:MsgSumLen])
		if verifySumCode != msg.Sum {
			return nil, errors.New("network data sum verify error"), nil
		}
		msg.Id = utils.ProtocalNumber(utils.GetProtoName(&outside.C2SGameCert{}))
		_CertProto := &outside.C2SGameCert{
			LoginSeq: utils.BytesToLong(binaryData[MsgSumLen+MsgTypeLen+MsgIsPos : MsgSumLen+MsgTypeLen+MsgIsPos+8]),
			Token:    binaryData[MsgSumLen+MsgTypeLen+MsgIsPos+8:],
		}
		content, _ := proto.Marshal(_CertProto)
		msg.data = content
		return msg, nil, nil
	case 3:
		msg := &LogicMessage{}
		if err := binary.Read(dataBuff, binary.BigEndian, &msg.Sum); err != nil {
			return nil, err, nil
		}
		msg.Sum = utils.BytesToLong(binaryData[:MsgSumLen])
		if verifySumCode != msg.Sum {
			return nil, errors.New("network data sum verify error"), nil
		}
		msg.MsgType = type_
		pGuid := binaryData[MsgSumLen+MsgTypeLen+MsgIsPos+4 : MsgSumLen+MsgTypeLen+MsgIsPos+4+8]
		pEnum := binaryData[MsgSumLen+MsgTypeLen+MsgIsPos+4+8 : MsgSumLen+MsgTypeLen+MsgIsPos+8+8]
		msg.PID = int64(utils.BytesToInt(pEnum))
		msg.data = binaryData[MsgSumLen+MsgTypeLen+MsgIsPos+8+8:]
		return msg, nil, func(conn network.IConn) {
			data := []byte{}
			data = append(append(data, byte(4)), pGuid...)
			msgData := append(append(utils.IntToBytes(8+1+8), utils.LongToBytes(int64(utils.GenerateCRCCheckCode(data)))...), data...)
			conn.SendByte(msgData)
		}
	}
	return nil, nil, nil
}
