package GateManager

import (
	"github.com/yaice-rx/yaice"
)

func Run(type_ string, groupId string, allowConn bool) {
	server := yaice.NewServer(type_, groupId, allowConn)
	server.AdaptationNetwork("http")
	server.Serve()
}
