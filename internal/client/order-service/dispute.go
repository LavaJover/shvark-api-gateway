package orderservice

import (
	"context"
	"time"

	orderpb "github.com/LavaJover/shvark-order-service/proto/gen/order"
	"google.golang.org/protobuf/types/known/durationpb"
)

func (c *OrderClient) CreateDispute(
	orderID, proofUrl, disputeReason string,
	ttl time.Duration,
	disputeAmountFiat float64,
) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	disputeResponse, err := c.orderService.CreateOrderDispute(
		ctx,
		&orderpb.CreateOrderDisputeRequest{
			OrderId: orderID,
			ProofUrl: proofUrl,
			DisputeReason: disputeReason,
			Ttl: durationpb.New(ttl),
			DisputeAmountFiat: disputeAmountFiat,
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

	_, err := c.orderService.AcceptOrderDispute(
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

	_, err := c.orderService.RejectOrderDispute(
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

	disputeResponse, err := c.orderService.GetOrderDisputeInfo(
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

	_, err := c.orderService.FreezeOrderDispute(
		ctx,
		&orderpb.FreezeOrderDisputeRequest{
			DisputeId: disputeID,
		},
	)

	return err
}

func (c *OrderClient) GetOrderDisputes(r *orderpb.GetOrderDisputesRequest) (*orderpb.GetOrderDisputesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.orderService.GetOrderDisputes(
		ctx,
		r,
	)
}