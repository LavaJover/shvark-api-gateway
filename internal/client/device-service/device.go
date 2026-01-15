package deviceservice

import (
	"context"
	"time"

	orderpb "github.com/LavaJover/shvark-order-service/proto/gen/order"
)

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