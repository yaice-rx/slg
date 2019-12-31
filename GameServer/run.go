package GameServer

import (
	"context"
	"github.com/yaice-rx/yaice"
	"google.golang.org/grpc"
	"slg/GameServer/rpc"
	"slg/Rpc"
)

func Run(type_ string, groupId string) {
	server := yaice.NewServer(type_, groupId, []string{"127.0.0.1:2379"})
	server.AddClientRpc(func(c *grpc.ClientConn) {
		client := auth_game_proto.NewVerifyUserServiceClient(c)
		client.VerifyUser(context.Background(), &auth_game_proto.A2GVerifyUser{PlayerGuid: groupId})
	})
	server.AddServerRpc(func(s *grpc.Server) {
		auth_game_proto.RegisterVerifyUserServiceServer(s, &rpc.VerifyUserService{})
	})
	server.MatchNetwork("tcp")
	server.Serve(20001, 20100)
}
