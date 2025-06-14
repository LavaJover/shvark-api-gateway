package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderpb "github.com/LavaJover/shvark-order-service/proto/gen"
)

type OrderClient struct {
	conn *grpc.ClientConn
	service orderpb.OrderServiceClient
}

func NewOrderClient(addr string) (*OrderClient, error) {
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

	return &OrderClient{
		conn: conn,
		service: orderpb.NewOrderServiceClient(conn),
	}, nil
}

func (c *OrderClient) CreateOrder(orderRequest *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.CreateOrder(
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

func (c *OrderClient) GetOrdersByTraderID(traderID string) (*orderpb.GetOrdersByTraderIDResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.GetOrdersByTraderID(
		ctx,
		&orderpb.GetOrdersByTraderIDRequest{
			TraderId: traderID,
		},
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

func (c *OrderClient) OpenOrderDispute(orderID string) (*orderpb.OpenOrderDisputeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.OpenOrderDispute(
		ctx,
		&orderpb.OpenOrderDisputeRequest{
			OrderId: orderID,
		},
	)
}

func (c *OrderClient) ResolveOrderDispute(orderID string) (*orderpb.ResolveOrderDisputeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.ResolveOrderDispute(
		ctx,
		&orderpb.ResolveOrderDisputeRequest{
			OrderId: orderID,
		},
	)
}