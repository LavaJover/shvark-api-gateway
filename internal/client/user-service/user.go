package userservice

import (
	"context"
	"time"

	"github.com/LavaJover/shvark-api-gateway/internal/domain"
	userpb "github.com/LavaJover/shvark-user-service/proto/gen"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

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

func (c *UserClient) UpdateUser(userID string, userData *domain.UpdateUserData, fields []string) (*userpb.UpdateUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.UpdateUser(
		ctx,
		&userpb.UpdateUserRequest{
			UserId: userID,
			User: &userpb.User{
				UserId:      userID,
				Login:       userData.Login,
				Username:    userData.Username,
				Password:    userData.Password,
				TwoFaSecret: userData.TwoFaSecret,
				Role:        userData.Role,
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

// GetUsersWithFilter получает пользователей с фильтрацией, сортировкой и пагинацией
func (c *UserClient) GetUsersWithFilter(r *userpb.GetUsersWithFilterRequest) (*userpb.GetUsersWithFilterResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.GetUsersWithFilter(ctx, r)
}

// DeleteUser удаляет пользователя
func (c *UserClient) DeleteUser(r *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.DeleteUser(ctx, r)
}

// ChangePassword изменяет пароль пользователя
func (c *UserClient) ChangePassword(r *userpb.ChangePasswordRequest) (*userpb.ChangePasswordResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.ChangePassword(ctx, r)
}

// SetTwoFaSecret устанавливает секрет 2FA
func (c *UserClient) SetTwoFaSecret(r *userpb.SetTwoFaSecretRequest) (*userpb.SetTwoFaSecretResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.SetTwoFaSecret(ctx, r)
}

// GetTwoFaSecretByID получает секрет 2FA по ID пользователя
func (c *UserClient) GetTwoFaSecretByID(r *userpb.GetTwoFaSecretByIDRequest) (*userpb.GetTwoFaSecretByIDResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.GetTwoFaSecretByID(ctx, r)
}

// GetUsers получает пользователей с пагинацией (базовая версия)
func (c *UserClient) GetUsers(r *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.GetUsers(ctx, r)
}

// GetTeamTraders получает трейдеров команды по ID тимлида
func (c *UserClient) GetTeamTraders(r *userpb.GetTeamTradersRequest) (*userpb.GetTeamTradersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.GetTeamTraders(ctx, r)
}

// AddTraderToTeam добавляет трейдера в команду
func (c *UserClient) AddTraderToTeam(r *userpb.AddTraderToTeamRequest) (*userpb.AddTraderToTeamResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.AddTraderToTeam(ctx, r)
}

// UpdateRelationshipParams обновляет параметры отношения
func (c *UserClient) UpdateRelationshipParams(r *userpb.UpdateRelationshipParamsRequest) (*userpb.UpdateRelationshipParamsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.UpdateRelationshipParams(ctx, r)
}

// GetRelationshipByID получает отношение по ID
func (c *UserClient) GetRelationshipByID(r *userpb.GetRelationshipByIDRequest) (*userpb.GetRelationshipByIDResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.GetRelationshipByID(ctx, r)
}