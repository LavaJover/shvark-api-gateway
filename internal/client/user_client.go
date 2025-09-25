package client

import (
	"context"
	"time"

	"github.com/LavaJover/shvark-api-gateway/internal/domain"
	userpb "github.com/LavaJover/shvark-user-service/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type UserClient struct {
	conn *grpc.ClientConn
	service userpb.UserServiceClient
}

func NewUserClient(addr string) (*UserClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

	return &UserClient{
		conn: conn,
		service: userpb.NewUserServiceClient(conn),
	}, nil
}

func (c *UserClient) CreateUser(login, username, password string) (*userpb.CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.CreateUser(
		ctx,
		&userpb.CreateUserRequest{
			Login: login,
			Username: username,
			Password: password,
		},
	)
}

func (c *UserClient) UpdateUser(userID string, user *domain.User, fields []string) (*userpb.UpdateUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.UpdateUser(
		ctx,
		&userpb.UpdateUserRequest{
			UserId: userID,
			User: &userpb.User{
				UserId: user.ID,
				Login: user.Login,
				Username: user.Username,
				Password: user.Password,
				TwoFaSecret: user.TwoFaSecret,
			},
			UpdateMask: &fieldmaskpb.FieldMask{Paths: fields},
		},
	)
}

func (c *UserClient) GetUserByID(userID string) (*userpb.GetUserByIDResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.GetUserByID(
		ctx,
		&userpb.GetUserByIDRequest{
			UserId: userID,
		},
	)
}

func (c *UserClient) GetUserByLogin(login string) (*userpb.GetUserByLoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.GetUserByLogin(
		ctx,
		&userpb.GetUserByLoginRequest{
			Login: login,
		},
	)
}

func (c *UserClient) GetTraders() (*userpb.GetTradersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.GetTraders(
		ctx,
		&userpb.GetTradersRequest{},
	)
}

func (c *UserClient) GetMerchants() (*userpb.GetMerchantsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.GetMerchants(
		ctx,
		&userpb.GetMerchantsRequest{},
	)
}

func (c *UserClient) PromoteToTeamLead(r *userpb.PromoteToTeamLeadRequest) (*userpb.PromoteToTeamLeadResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.PromoteToTeamLead(
		ctx,
		r,
	)
}

func (c *UserClient) DemoteTeamLead(r *userpb.DemoteTeamLeadRequest) (*userpb.DemoteTeamLeadResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.DemoteTeamLead(
		ctx,
		r,
	)
}

func (c *UserClient) GetUsersByRole(r *userpb.GetUsersByRoleRequest) (*userpb.GetUsersByRoleResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.GetUsersByRole(
		ctx,
		r,
	)
}

func (c *UserClient) SetTwoFaEnabled(r *userpb.SetTwoFaEnabledRequest) (*userpb.SetTwoFaEnabledResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.SetTwoFaEnabled(
		ctx,
		r,
	)
}