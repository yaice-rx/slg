package Service

type Token struct {
	Port      int    `json:"port"`
	Result    int    `json:"result"`
	Guid      uint64 `json:"guid"`
	SessionId uint64 `json:"sessionId"`
	Host      string `json:"host"`
}
