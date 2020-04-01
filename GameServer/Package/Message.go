package Package

type AuthMessage struct {
	MsgType  uint8 //消息类型
	Sum      int64 //消息的ID
	IsPos    int64 //位置编号
	LoginSeq int64 //登陆编码
	Id       int32 //消息编码
	data     []byte
}

func (m *AuthMessage) GetMsgId() int32 {
	return m.Id
}

func (m *AuthMessage) GetData() []byte {
	return m.data
}

type LogicMessage struct {
	MsgType uint8 //消息类型
	MsgLen  uint32
	Sum     int64 //消息的ID
	IsPos   int64 //位置编号
	PID     int64
	PEnum   int32
	data    []byte
}

func (m *LogicMessage) GetMsgId() int32 {
	return m.PEnum
}

func (m *LogicMessage) GetData() []byte {
	return m.data
}
