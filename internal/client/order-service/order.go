package orderservice

import (
	"context"
	"time"

	orderpb "github.com/LavaJover/shvark-order-service/proto/gen/order"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *OrderClient) CreatePayInOrder(orderRequest *orderpb.CreatePayInOrderRequest) (*orderpb.CreatePayInOrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.CreatePayInOrder(
		ctx,
		orderRequest,
	)
}

func (c *OrderClient) CreatePayOutOrder(orderRequest *orderpb.CreatePayOutOrderRequest) (*orderpb.CreatePayOutOrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.CreatePayOutOrder(
		ctx,
		orderRequest,
	)
}

func (c *OrderClient) GetOrderByID(orderID string) (*orderpb.GetOrderByIDResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.GetOrderByID(
		ctx,
		&orderpb.GetOrderByIDRequest{
			OrderId: orderID,
		},
	)
}

func (c *OrderClient) GetOrderByMerchantOrderID(merchantOrderID string) (*orderpb.GetOrderByMerchantOrderIDResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.GetOrderByMerchantOrderID(
		ctx,
		&orderpb.GetOrderByMerchantOrderIDRequest{
			MerchantOrderId: merchantOrderID,
		},
	)
}

func (c *OrderClient) GetOrdersByTraderID(request *orderpb.GetOrdersByTraderIDRequest) (*orderpb.GetOrdersByTraderIDResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.GetOrdersByTraderID(
		ctx,
		request,
	)
}

func (c *OrderClient) ApproveOrder(orderID string) (*orderpb.ApproveOrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.ApproveOrder(
		ctx,
		&orderpb.ApproveOrderRequest{
			OrderId: orderID,
		},
	)
}

func (c *OrderClient) CancelOrder(orderID string) (*orderpb.CancelOrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.CancelOrder(
		ctx,
		&orderpb.CancelOrderRequest{
			OrderId: orderID,
		},
	)
}

type OrderStats struct {
	TotalOrders 			int64 	
	SucceedOrders 			int64 	
	CanceledOrders 			int64 	
	ProcessedAmountFiat 	float64 
	ProcessedAmountCrypto 	float64 
	CanceledAmountFiat 		float64 
	CanceledAmountCrypto 	float64 
	IncomeCrypto 			float64 
}

func (c *OrderClient) GetOrderStats(
	traderID string,
	dateFrom, dateTo time.Time,
) (*orderpb.GetOrderStatisticsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.GetOrderStatistics(
		ctx,
		&orderpb.GetOrderStatisticsRequest{
			TraderId: traderID,
			DateFrom: timestamppb.New(dateFrom),
			DateTo: timestamppb.New(dateTo),
		},
	)
}

func (c *OrderClient) GetOrders(r *orderpb.GetOrdersRequest) (*orderpb.GetOrdersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.service.GetOrders(
		ctx,
		r,
	)
}

func (c *OrderClient) GetAllOrders(r *orderpb.GetAllOrdersRequest) (*orderpb.GetAllOrdersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.service.GetAllOrders(
		ctx,
		r,
	)
}