package Session

import (
	"SLGGAME/Service"
)

//auth在线服务器连接
var AuthContainsGameMgr = Service.NewConnSessionMgr()

//玩家在线列表
var PlayerContainsGameMgr = Service.NewConnSessionMgr()
