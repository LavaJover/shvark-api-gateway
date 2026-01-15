package ssoservice

import (
	"context"
	"time"

	ssopb "github.com/LavaJover/shvark-sso-service/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SSOClient struct {
	conn *grpc.ClientConn
	service ssopb.SSOServiceClient
}

func NewSSOClient(addr string) (*SSOClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
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

	return &SSOClient{
		conn: conn,
		service: ssopb.NewSSOServiceClient(conn),
	}, err
}