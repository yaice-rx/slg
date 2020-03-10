package Package

type Message struct {
	DataLen uint32 //消息的长度
	Sum     int64  //消息的ID
	MsgType uint16 //消息类型
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) GetMsgId() int32 {
	return 0
}
