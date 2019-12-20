package GameServer

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/yaice-rx/yaice"
	"github.com/yaice-rx/yaice/cron"
	"github.com/yaice-rx/yaice/network"
	proto_ "github.com/yaice-rx/yaice/proto"
)

func Run(type_ string, groupId string, allowConn bool) {
	server := yaice.NewServer(type_, groupId, allowConn)
	server.AddRouter(&proto_.S2CServiceAssociate{}, func(conn network.IConn, content []byte) {
		var data proto_.S2CServiceAssociate
		if json.Unmarshal(content, &data) != nil {
			return
		}
		logrus.Println("接收到服务器的消息数据", data.Pid, data.TypeName)
		//心跳
		cron.CronMgr.AddCronTask(10, -1, func() {
			data := proto_.C2SServicePing{}
			if err := conn.SendMsg(&data); err != nil {
				logrus.Error("send err ", err)
				return
			}
		})
	})
	server.AdaptationNetwork("tcp")
	server.Serve()
}
