package orderservice

import (
	"context"
	"time"

	orderpb "github.com/LavaJover/shvark-order-service/proto/gen/order"
)

func (c *OrderClient) CreateTeamRelation(r *orderpb.CreateTeamRelationRequest) (*orderpb.CreateTeamRelationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.teamRelationsService.CreateTeamRelation(
		ctx,
		r,
	)
}

func (c *OrderClient) UpdateTeamRelationParams(r *orderpb.UpdateRelationParamsRequest) (*orderpb.UpdateRelationParamsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.teamRelationsService.UpdateRelationParams(
		ctx,
		r,
	)
}

func (c *OrderClient) GetTeamRelationsByTeamLeadID(r *orderpb.GetRelationsByTeamLeadIDRequest) (*orderpb.GetRelationsByTeamLeadIDResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.teamRelationsService.GetRelationsByTeamLeadID(
		ctx,
		r,
	)
}

func (c *OrderClient) DeleteTeamRelationship(r *orderpb.DeleteTeamRelationshipRequest) (*orderpb.DeleteTeamRelationshipResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.teamRelationsService.DeleteTeamRelationship(
		ctx,
		r,
	)
}