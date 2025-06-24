package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"

	orderpb "github.com/LavaJover/shvark-order-service/proto/gen"
)

type OrderClient struct {
	conn *grpc.ClientConn
	service orderpb.OrderServiceClient
	trafficService orderpb.TrafficServiceClient
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
		trafficService: orderpb.NewTrafficServiceClient(conn),
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

func (c *OrderClient) AddTraffic(
	merchantID, traderID string,
	traderReward, traderPriority float64,
	enabled bool,
	) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.trafficService.AddTraffic(
		ctx,
		&orderpb.AddTrafficRequest{
			MerchantId: merchantID,
			TraderId: traderID,
			TraderRewardPercent: traderReward,
			TraderPriority: traderPriority,
			Enabled: enabled,
		},
	)
	return err
}

func (c *OrderClient) EditTraffic(
	trafficID string,
	traderReward, traderPriority float64,
	enabled bool,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.trafficService.EditTraffic(
		ctx,
		&orderpb.EditTrafficRequest{
			Traffic: &orderpb.Traffic{
				Id: trafficID,
				TraderRewardPercent: traderReward,
				TraderPriority: traderPriority,
				Enabled: enabled,
			},
		},
	)

	return err
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

func (c *OrderClient) CreateDispute(
	orderID, proofUrl, disputeReason string,
	ttl time.Duration,
) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	disputeResponse, err := c.service.CreateOrderDispute(
		ctx,
		&orderpb.CreateOrderDisputeRequest{
			OrderId: orderID,
			ProofUrl: proofUrl,
			DisputeReason: disputeReason,
			Ttl: durationpb.New(ttl),
		},
	)

	if err != nil {
		return "", err
	}

	return disputeResponse.DisputeId, nil
}

func (c *OrderClient) AcceptDispute(
	disputeID string,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.service.AcceptOrderDispute(
		ctx,
		&orderpb.AcceptOrderDisputeRequest{
			DisputeId: disputeID,
		},
	)

	return err
}

func (c *OrderClient) RejectDispute(
	disputeID string,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.service.RejectOrderDispute(
		ctx,
		&orderpb.RejectOrderDisputeRequest{
			DisputeId: disputeID,
		},
	)
	return err
}

type Dispute struct {
	DisputeID 	  string
	OrderID 	  string
	ProofUrl 	  string
	DisputeReason string
	DisputeStatus string
}

func (c *OrderClient) GetDisputeInfo(
	disputeID string,
) (*Dispute, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	disputeResponse, err := c.service.GetOrderDisputeInfo(
		ctx,
		&orderpb.GetOrderDisputeInfoRequest{
			DisputeId: disputeID,
		},
	)

	if err != nil {
		return nil, err
	}

	return &Dispute{
		DisputeID: disputeResponse.Dispute.DisputeId,
		OrderID: disputeResponse.Dispute.OrderId,
		ProofUrl: disputeResponse.Dispute.ProofUrl,
		DisputeReason: disputeResponse.Dispute.DisputeReason,
		DisputeStatus: disputeResponse.Dispute.DisputeStatus,
	}, nil
}

func (c *OrderClient) FreeezeDispute(disputeID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.service.FreezeOrderDispute(
		ctx,
		&orderpb.FreezeOrderDisputeRequest{
			DisputeId: disputeID,
		},
	)

	return err
}