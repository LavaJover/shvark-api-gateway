package client

import (
	"context"
	"time"

	"github.com/LavaJover/shvark-api-gateway/internal/domain"
	bankingpb "github.com/LavaJover/shvark-banking-service/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
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

func (c *BankingClient) CreateBankDetail(bankDetail *domain.BankDetail) (*bankingpb.CreateBankDetailResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.CreateBankDetail(
		ctx,
		&bankingpb.CreateBankDetailRequest{
			TraderId: bankDetail.TraderID,
			Currency: bankDetail.Currency,
			Country: bankDetail.Country,
			MinAmount: float64(bankDetail.MinAmount),
			MaxAmount: float64(bankDetail.MaxAmount),
			BankName: bankDetail.BankName,
			PaymentSystem: bankDetail.PaymentSystem,
			Enabled: bankDetail.Enabled,
			Delay: durationpb.New(bankDetail.Delay),
			CardNumber: bankDetail.CardNumber,
			Phone: bankDetail.Phone,
			Owner: bankDetail.Owner,
			MaxOrdersSimultaneosly: bankDetail.MaxOrdersSimultaneosly,
			MaxAmountDay: float64(bankDetail.MaxAmountDay),
			MaxAmountMonth: float64(bankDetail.MaxAmountMonth),
		},
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

func (c *BankingClient) UpdateBankDetail(bankDetail *domain.BankDetail) (*bankingpb.UpdateBankDetailResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.UpdateBankDetail(
		ctx,
		&bankingpb.UpdateBankDetailRequest{
			BankDetail: &bankingpb.BankDetail{
				BankDetailId: bankDetail.ID,
				TraderId: bankDetail.TraderID,
				Currency: bankDetail.Currency,
				Country: bankDetail.Country,
				MinAmount: float64(bankDetail.MinAmount),
				MaxAmount: float64(bankDetail.MaxAmount),
				BankName: bankDetail.BankName,
				PaymentSystem: bankDetail.PaymentSystem,
				Enabled: bankDetail.Enabled,
				Delay: durationpb.New(bankDetail.Delay),
				CardNumber: bankDetail.CardNumber,
				Phone: bankDetail.Phone,
				Owner: bankDetail.Owner,
				MaxOrdersSimultaneosly: bankDetail.MaxOrdersSimultaneosly,
				MaxAmountDay: float64(bankDetail.MaxAmountDay),
				MaxAmountMonth: float64(bankDetail.MaxAmountMonth),
			},
		},
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