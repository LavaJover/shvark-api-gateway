package orderservice

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	orderpb "github.com/LavaJover/shvark-order-service/proto/gen/order"
)

func (c *OrderClient) CreateStore(r *orderpb.CreateStoreRequest) (*orderpb.CreateStoreResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.storeService.CreateStore(ctx, r)
}

func (c *OrderClient) GetStore(r *orderpb.GetStoreRequest) (*orderpb.GetStoreResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.storeService.GetStore(ctx, r)
}

func (c *OrderClient) UpdateStore(r *orderpb.UpdateStoreRequest) (*orderpb.UpdateStoreResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.storeService.UpdateStore(ctx, r)
}

func (c *OrderClient) DeleteStore(r *orderpb.DeleteStoreRequest) (*orderpb.DeleteStoreResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.storeService.DeleteStore(ctx, r)
}

func (c *OrderClient) ListStores(r *orderpb.ListStoresRequest) (*orderpb.ListStoresResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.storeService.ListStores(ctx, r)
}

func (c *OrderClient) GetStoreWithTraffics(r *orderpb.GetStoreWithTrafficsRequest) (*orderpb.GetStoreWithTrafficsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.storeService.GetStoreWithTraffics(ctx, r)
}

func (c *OrderClient) GetStoreByTrafficId(r *orderpb.GetStoreByTrafficIdRequest) (*orderpb.GetStoreByTrafficIdResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.storeService.GetStoreByTrafficId(ctx, r)
}

func (c *OrderClient) CheckStoreNameUnique(r *orderpb.CheckStoreNameUniqueRequest) (*orderpb.CheckStoreNameUniqueResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.storeService.CheckStoreNameUnique(ctx, r)
}

func (c *OrderClient) ValidateStoreForTraffic(r *orderpb.ValidateStoreForTrafficRequest) (*orderpb.ValidateStoreForTrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.storeService.ValidateStoreForTraffic(ctx, r)
}

func (c *OrderClient) ToggleStoreStatus(r *orderpb.ToggleStoreStatusRequest) (*orderpb.ToggleStoreStatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.storeService.ToggleStoreStatus(ctx, r)
}

func (c *OrderClient) EnableStore(r *orderpb.EnableStoreRequest) (*orderpb.EnableStoreResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.storeService.EnableStore(ctx, r)
}

func (c *OrderClient) DisableStore(r *orderpb.DisableStoreRequest) (*orderpb.DisableStoreResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.storeService.DisableStore(ctx, r)
}

func (c *OrderClient) BulkUpdateStoresStatus(r *orderpb.BulkUpdateStoresStatusRequest) (*orderpb.BulkUpdateStoresStatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.storeService.BulkUpdateStoresStatus(ctx, r)
}

func (c *OrderClient) SearchStores(r *orderpb.SearchStoresRequest) (*orderpb.SearchStoresResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.storeService.SearchStores(ctx, r)
}

func (c *OrderClient) GetStoresByMerchant(r *orderpb.GetStoresByMerchantRequest) (*orderpb.GetStoresByMerchantResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.storeService.GetStoresByMerchant(ctx, r)
}

func (c *OrderClient) GetActiveStores(r *orderpb.GetActiveStoresRequest) (*orderpb.GetActiveStoresResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.storeService.GetActiveStores(ctx, r)
}

func (c *OrderClient) GetStoreMetrics(r *orderpb.GetStoreMetricsRequest) (*orderpb.GetStoreMetricsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.storeService.GetStoreMetrics(ctx, r)
}

func (c *OrderClient) CalculateStoreMetrics(r *orderpb.CalculateStoreMetricsRequest) (*orderpb.CalculateStoreMetricsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.storeService.CalculateStoreMetrics(ctx, r)
}

func (c *OrderClient) BatchGetStores(r *orderpb.BatchGetStoresRequest) (*orderpb.BatchGetStoresResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.storeService.BatchGetStores(ctx, r)
}

func (c *OrderClient) HealthCheck() (*orderpb.HealthCheckResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.storeService.HealthCheck(ctx, &emptypb.Empty{})
}