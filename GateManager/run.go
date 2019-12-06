package GateManager

import (
	"fmt"
	"github.com/yaice-rx/yaice"
	"net/http"
)

func Run(type_ string, groupId string, allowConn bool) {
	server := yaice.NewServer(type_, groupId, allowConn)
	server.AdaptationNetwork("http")
	server.AddHttpHandler("/login", func(write http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		fmt.Println(request.Form)
	})
	server.Serve()
}
