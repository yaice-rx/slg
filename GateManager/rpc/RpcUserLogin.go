package rpc

import (
	"context"
	"fmt"
	"slg/Rpc"
)

type VerifyUserService struct {
	auth_game_proto.UnimplementedVerifyUserServiceServer
}

func (u *VerifyUserService) VerifyUser(ctx context.Context, point *auth_game_proto.A2GVerifyUser) (*auth_game_proto.G2CVerifyUser, error) {
	fmt.Println("grpc data content =>", point)
	return &auth_game_proto.G2CVerifyUser{LoginState: true}, nil
}
