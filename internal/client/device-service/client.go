package deviceservice

import (
	"context"
	"time"

	orderpb "github.com/LavaJover/shvark-order-service/proto/gen/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DeviceClient struct {
	conn *grpc.ClientConn
    client orderpb.DeviceServiceClient
}

func NewDeviceClient(addr string) (*DeviceClient, error) {
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

	return &DeviceClient{
		conn: conn,
		client: orderpb.NewDeviceServiceClient(conn),
	}, nil
}