package orderservice

import (
	"context"
	"time"

	orderpb "github.com/LavaJover/shvark-order-service/proto/gen/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderClient struct {
	conn *grpc.ClientConn
	orderService orderpb.OrderServiceClient
	trafficService orderpb.TrafficServiceClient
	bankDetailService orderpb.BankDetailServiceClient
	teamRelationsService orderpb.TeamRelationsServiceClient
	deviceService orderpb.DeviceServiceClient
	antifraudService orderpb.AntiFraudServiceClient
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
		orderService: orderpb.NewOrderServiceClient(conn),
		trafficService: orderpb.NewTrafficServiceClient(conn),
		bankDetailService: orderpb.NewBankDetailServiceClient(conn),
		teamRelationsService: orderpb.NewTeamRelationsServiceClient(conn),
		deviceService: orderpb.NewDeviceServiceClient(conn),
		antifraudService: orderpb.NewAntiFraudServiceClient(conn),
	}, nil
}