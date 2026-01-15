package orderservice

import (
	"context"
	"time"

	orderpb "github.com/LavaJover/shvark-order-service/proto/gen/order"
)

// ============= СТАТУС УСТРОЙСТВ =============

// UpdateDeviceLiveness обновляет статус онлайн устройства
func (c *OrderClient) UpdateDeviceLiveness(ctx context.Context, req *orderpb.UpdateDeviceLivenessRequest) (*orderpb.UpdateDeviceLivenessResponse, error) {
    timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    return c.deviceService.UpdateDeviceLiveness(timeoutCtx, req)
}

// GetDeviceStatus получает статус устройства
func (c *OrderClient) GetDeviceStatus(ctx context.Context, req *orderpb.GetDeviceStatusRequest) (*orderpb.GetDeviceStatusResponse, error) {
    timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    return c.deviceService.GetDeviceStatus(timeoutCtx, req)
}

// GetTraderDevicesStatus получает статусы всех устройств трейдера
func (c *OrderClient) GetTraderDevicesStatus(ctx context.Context, req *orderpb.GetTraderDevicesStatusRequest) (*orderpb.GetTraderDevicesStatusResponse, error) {
    timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    return c.deviceService.GetTraderDevicesStatus(timeoutCtx, req)
}

func (c *OrderClient) GetTraderDevices(r *orderpb.GetTraderDevicesRequest) (*orderpb.GetTraderDevicesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.deviceService.GetTraderDevices(
		ctx,
		r,
	)
}

func (c *OrderClient) CreateDevice(r *orderpb.CreateDeviceRequest) (*orderpb.CreateDeviceResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.deviceService.CreateDevice(
		ctx,
		r,
	)
}

func (c *OrderClient) DeleteDevice(r *orderpb.DeleteDeviceRequest) (*orderpb.DeleteDeviceResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.deviceService.DeleteDevice(
		ctx,
		r,
	)
}

func (c *OrderClient) EditeDevice(r *orderpb.EditDeviceRequest) (*orderpb.EditDeviceResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.deviceService.EditDevice(
		ctx,
		r,
	)
}