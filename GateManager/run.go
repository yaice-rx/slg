package GateManager

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/yaice-rx/yaice"
	"github.com/yaice-rx/yaice/network"
	proto_ "github.com/yaice-rx/yaice/proto"
	"time"
)

func Run(type_ string, groupId string, allowConn bool) {
	server := yaice.NewServer(type_, groupId, allowConn)
	server.AddRouter(&proto_.C2SServiceAssociate{}, func(conn network.IConn, content []byte) {
		var data proto_.C2SServiceAssociate
		if err := json.Unmarshal(content, &data); err != nil {
			logrus.Println(""+err.Error(), "====", string(content))
			return
		}
		logrus.Println(data.Pid, data.TypeName)
		conn.SendMsg(&proto_.S2CServiceAssociate{
			TypeName: type_,
			Pid:      123,
		})
	})
	server.AddRouter(&proto_.C2SServicePing{}, func(conn network.IConn, content []byte) {
		logrus.Println("ping =================== ", time.Now().String())
	})
	server.AdaptationNetwork("http")
	server.Serve()
}
