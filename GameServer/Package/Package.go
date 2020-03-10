package Package

import (
	"github.com/yaice-rx/yaice/network"
)

type Package struct {
}

const (
	MsgLengthLen = 4
	MsgSumLen    = 8
	MsgTypeLen   = 1
)

func NewPackage() network.IPacket {
	return &Package{}
}

func (p *Package) GetHeadLen() uint32 {
	return MsgLengthLen + MsgSumLen + MsgTypeLen
}

func (p *Package) Pack(msg network.TransitData) []byte {
	return nil
}

func (p *Package) Unpack(binaryData []byte) (network.IMessage, error) {

	return nil, nil
}
