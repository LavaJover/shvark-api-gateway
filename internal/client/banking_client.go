package client

import (
	"context"
	"time"

	bankingpb "github.com/LavaJover/shvark-banking-service/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BankingClient struct {
	conn *grpc.ClientConn
	service bankingpb.BankingServiceClient
}

func NewBankingClient(addr string) (*BankingClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
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
	return &BankingClient{
		conn: conn,
		service: bankingpb.NewBankingServiceClient(conn),
	}, nil
}

func (c *BankingClient) CreateBankDetail(bankDetailRequest *bankingpb.CreateBankDetailRequest) (*bankingpb.CreateBankDetailResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.CreateBankDetail(
		ctx,
		bankDetailRequest,
	)
}

func (c *BankingClient) DeleteBankDetail(bankDetailID string) (*bankingpb.DeleteBankDetailResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.DeleteBankDetail(
		ctx,
		&bankingpb.DeleteBankDetailRequest{
			BankDetailId: bankDetailID,
		},
	)
}

func (c *BankingClient) GetBankDetailByID(bankDetailID string) (*bankingpb.GetBankDetailByIDResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.GetBankDetailByID(
		ctx,
		&bankingpb.GetBankDetailByIDRequest{
			BankDetailId: bankDetailID,
		},
	)
}

func (c *BankingClient) UpdateBankDetail(bankDetailRequest *bankingpb.UpdateBankDetailRequest) (*bankingpb.UpdateBankDetailResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.UpdateBankDetail(
		ctx,
		bankDetailRequest,
	)
}

func (c *BankingClient) GetBankDetailsByTraderID(traderID string) (*bankingpb.GetBankDetailsByTraderIDResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.GetBankDetailsByTraderID(
		ctx,
		&bankingpb.GetBankDetailsByTraderIDRequest{
			TraderId: traderID,
		},
	)
}