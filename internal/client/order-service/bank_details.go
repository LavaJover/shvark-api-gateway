package orderservice

import (
	"context"
	"time"

	orderpb "github.com/LavaJover/shvark-order-service/proto/gen/order"
)

func (c *OrderClient) CreateBankDetail(createBankDetailRequest *orderpb.CreateBankDetailRequest) (*orderpb.CreateBankDetailResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.bankDetailService.CreateBankDetail(
		ctx,
		createBankDetailRequest,
	)
}

func (c *OrderClient) EditBankDetail(editBankDetailRequest *orderpb.UpdateBankDetailRequest) (*orderpb.UpdateBankDetailResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.bankDetailService.UpdateBankDetail(
		ctx,
		editBankDetailRequest,
	)
}

func (c *OrderClient) DeleteBankDetail(deleteBankDetailRequest *orderpb.DeleteBankDetailRequest) (*orderpb.DeleteBankDetailResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.bankDetailService.DeleteBankDetail(
		ctx,
		deleteBankDetailRequest,
	)
}

func (c *OrderClient) GetBankDetailsByTraderID(getBankDetailsRequest *orderpb.GetBankDetailsByTraderIDRequest) (*orderpb.GetBankDetailsByTraderIDResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.bankDetailService.GetBankDetailsByTraderID(
		ctx,
		getBankDetailsRequest,
	)
}

func (c *OrderClient) GetBankDetailByID(getbankDetailRequest *orderpb.GetBankDetailByIDRequest) (*orderpb.GetBankDetailByIDResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.bankDetailService.GetBankDetailByID(
		ctx,
		getbankDetailRequest,
	)
}

func (c *OrderClient) GetBankDetailsStatsByTraderID(getStatsRequest *orderpb.GetBankDetailsStatsByTraderIDRequest) (*orderpb.GetBankDetailsStatsByTraderIDResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.bankDetailService.GetBankDetailsStatsByTraderID(
		ctx,
		getStatsRequest,
	)
}

func (c *OrderClient) GetBankDetails(r *orderpb.GetBankDetailsRequest) (*orderpb.GetBankDetailsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.bankDetailService.GetBankDetails(
		ctx,
		r,
	)
}