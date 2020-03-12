package Package

type AuthMessage struct {
	Sum      int64 //消息的ID
	MsgType  uint8 //消息类型
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
	Sum     int64 //消息的ID
	MsgType uint8 //消息类型
	IsPos   int64 //位置编号
	MsgLen  uint32
	PID     int64
	PEnum   int64
	data    []byte
}

func (m *LogicMessage) GetMsgId() int32 {
	return int32(m.PID)
}

func (m *LogicMessage) GetData() []byte {
	return m.data
}
