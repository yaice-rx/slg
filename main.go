package main

import (
	"flag"
	"fmt"
	"slg/GameServer"
	"slg/GateManager"
)

var Type = flag.String("type", "auth", "Input Server Type")
var Group = flag.String("group", "king_war", "Input Server Type")
var AllowConnect = flag.Bool("allowConnect", false, "")

func main() {
	flag.Parse()
	switch *Type {
	case "auth":
		GateManager.Run(*Type, *Group, *AllowConnect)
		break
	case "game":
		GameServer.Run(*Type, *Group, *AllowConnect)
		break
	default:
		fmt.Println("please select service type")
		break
	}
}
