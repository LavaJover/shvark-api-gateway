package client

import (
	"context"
	"time"

	profilepb "github.com/LavaJover/shvark-profile-service/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProfileClient struct {
	conn *grpc.ClientConn
	service profilepb.ProfileServiceClient
}

func NewProfileClient(addr string) (*ProfileClient, error) {
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

	return &ProfileClient{
		conn: conn,
		service: profilepb.NewProfileServiceClient(conn),
	}, nil
}

func (c *ProfileClient) GetProfileByID(profileID string) (*profilepb.GetProfileByIDResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.GetProfileByID(
		ctx, 
		&profilepb.GetProfileByIDRequest{
			ProfileId: profileID,
		},
	)
}