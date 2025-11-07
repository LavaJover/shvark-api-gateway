package client

import (
    "context"
    "time"

    orderpb "github.com/LavaJover/shvark-order-service/proto/gen"
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

// GetTraderDevicesStatus получает статусы всех устройств трейдера
func (c *DeviceClient) GetTraderDevicesStatus(ctx context.Context, req *orderpb.GetTraderDevicesStatusRequest) (*orderpb.GetTraderDevicesStatusResponse, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    return c.client.GetTraderDevicesStatus(ctx, req)
}

// GetDeviceStatus получает статус конкретного устройства
func (c *DeviceClient) GetDeviceStatus(ctx context.Context, req *orderpb.GetDeviceStatusRequest) (*orderpb.GetDeviceStatusResponse, error) {
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    return c.client.GetDeviceStatus(ctx, req)
}