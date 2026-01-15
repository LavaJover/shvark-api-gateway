package authzservice

import (
	"context"
	"time"

	authzpb "github.com/LavaJover/shvark-authz-service/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthzClient struct {
	conn *grpc.ClientConn
	service authzpb.AuthzServiceClient
}

func NewAuthzClient(addr string) (*AuthzClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		return nil, err
	}

	return &AuthzClient{
		conn: conn,
		service: authzpb.NewAuthzServiceClient(conn),
	}, nil
}