package client

import (
	"context"
	"errors"
	"time"

	ssopb "github.com/LavaJover/shvark-sso-service/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SSOClient struct {
	conn *grpc.ClientConn
	service ssopb.SSOServiceClient
}

func NewSSOClient(addr string) (*SSOClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
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

	return &SSOClient{
		conn: conn,
		service: ssopb.NewSSOServiceClient(conn),
	}, err
}

func (c *SSOClient) Register(login, username, rawPassword, role string) (*ssopb.RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.service.Register(ctx, &ssopb.RegisterRequest{
		Login: login,
		Username: username,
		Password: rawPassword,
		Role: role,
	})
}

func (c *SSOClient) RegisterWithretry(login, username, rawPassword, role string, maxRetries int) (*ssopb.RegisterResponse, error) {
	for range maxRetries {
		resp, err := c.Register(login, username, rawPassword, role)
		if err == nil {
			return resp, err
		}
		time.Sleep(1 * time.Second)
	}
	return nil, errors.New("max retries exceeded")
}

func (c *SSOClient) Login(login, rawPassword string, code string) (*ssopb.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.service.Login(ctx, &ssopb.LoginRequest{
		Login: login,
		Password: rawPassword,
		TwoFaCode: code,
	})
}

func (c *SSOClient) ValidateToken(token string) (*ssopb.ValidateTokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.service.ValidateToken(ctx, &ssopb.ValidateTokenRequest{
		AccessToken: token,
	})
}

func (c *SSOClient) GetUserByToken(token string) (*ssopb.GetUserByTokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.service.GetUserByToken(ctx, &ssopb.GetUserByTokenRequest{
		AccessToken: token,
	})
}

func (c *SSOClient) Setup2FA(userID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	response, err := c.service.Setup2FA(
		ctx,
		&ssopb.Setup2FARequest{
			UserId: userID,
		},
	)
	if err != nil {
		return "", err
	}

	return response.QrUrl, nil
}

func (c *SSOClient) Verify2FA(userID, code string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	response, err := c.service.Verify2FA(
		ctx,
		&ssopb.Verify2FARequest{
			UserId: userID,
			Code: code,
		},
	)
	if err != nil {
		return false, err
	}
	return response.Verif, nil
}

func (c *SSOClient) Close() {
	c.conn.Close()
}