package orderservice

import (
	"context"
	"time"

	orderpb "github.com/LavaJover/shvark-order-service/proto/gen/order"
)

// GetTraderTraffic получает все записи трафика для трейдера
func (c *OrderClient) GetTraderTraffic(r *orderpb.GetTraderTrafficRequest) (*orderpb.GetTraderTrafficResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return c.trafficService.GetTraderTraffic(ctx, r)
}

// GetTrafficLockStatuses получает статусы блокировки трафика
func (c *OrderClient) GetTrafficLockStatuses(r *orderpb.GetTrafficLockStatusesRequest) (*orderpb.GetTrafficLockStatusesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.trafficService.GetTrafficLockStatuses(ctx, r)
}

// CheckTrafficUnlocked проверяет, разблокирован ли трафик
func (c *OrderClient) CheckTrafficUnlocked(r *orderpb.CheckTrafficUnlockedRequest) (*orderpb.CheckTrafficUnlockedResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.trafficService.CheckTrafficUnlocked(ctx, r)
}

func (c *OrderClient) AddTraffic(r *orderpb.AddTrafficRequest) (*orderpb.AddTrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	response, err := c.trafficService.AddTraffic(
		ctx,
		r,
	)
	return response, err
}

func (c *OrderClient) EditTraffic(r *orderpb.EditTrafficRequest) (*orderpb.EditTrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	response, err := c.trafficService.EditTraffic(
		ctx,
		r,
	)

	return response, err
}

func (c *OrderClient) GetTrafficRecords(page, limit int32) ([]*orderpb.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	trafficResponse, err := c.trafficService.GetTrafficRecords(
		ctx,
		&orderpb.GetTrafficRecordsRequest{
			Page: page,
			Limit: limit,
		},
	)
	if err != nil {
		return nil, err
	}

	return trafficResponse.TrafficRecords, nil
}

func (c *OrderClient) DeleteTraffic(trafficID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.trafficService.DeleteTraffic(
		ctx,
		&orderpb.DeleteTrafficRequest{
			TrafficId: trafficID,
		},
	)

	return err
}

func (c *OrderClient) SetTraderLockTrafficStatus(r *orderpb.SetTraderLockTrafficStatusRequest) (*orderpb.SetTraderLockTrafficStatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.trafficService.SetTraderLockTrafficStatus(
		ctx,
		r,
	)
}

func (c *OrderClient) SetMerchantLockTrafficStatus(r *orderpb.SetMerchantLockTrafficStatusRequest) (*orderpb.SetMerchantLockTrafficStatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.trafficService.SetMerchantLockTrafficStatus(
		ctx,
		r,
	)
}

func (c *OrderClient) SetManuallyLockTrafficStatus(r *orderpb.SetManuallyLockTrafficStatusRequest) (*orderpb.SetManuallyLockTrafficStatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.trafficService.SetManuallyLockTrafficStatus(
		ctx,
		r,
	)
}

func (c *OrderClient) SetAntifraudLockTrafficStatus(r *orderpb.SetAntifraudLockTrafficStatusRequest) (*orderpb.SetAntifraudLockTrafficStatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.trafficService.SetAntifraudLockTrafficStatus(
		ctx,
		r,
	)
}