package authzservice

import (
	"context"
	"time"

	authzpb "github.com/LavaJover/shvark-authz-service/proto/gen"
)

func (c *AuthzClient) AssignRole(userID, role string) (*authzpb.AssignRoleResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.AssignRole(
		ctx,
		&authzpb.AssignRoleRequest{
			UserId: userID,
			Role: role,
		},
	)
}

func (c *AuthzClient) RevokeRole(userID, role string) (*authzpb.RevokeRoleResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.RevokeRole(
		ctx,
		&authzpb.RevokeRoleRequest{
			UserId: userID,
			Role: role,
		},
	)
}

func (c *AuthzClient) AddPolicy(role, object, action string) (*authzpb.AddPolicyResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.AddPolicy(
		ctx,
		&authzpb.AddPolicyRequest{
			Role: role,
			Object: object,
			Action: action,
		},
	)
}

func (c *AuthzClient) DeletePolicy(role, object, action string) (*authzpb.DeletePolicyResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.DeletePolicy(
		ctx,
		&authzpb.DeletePolicyRequest{
			Role: role,
			Object: object,
			Action: action,
		},
	)
}

func (c *AuthzClient) CheckPermission(userID, object, action string) (*authzpb.CheckPermissionResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.service.CheckPermission(
		ctx,
		&authzpb.CheckPermissionRequest{
			UserId: userID,
			Object: object,
			Action: action,
		},
	)
}