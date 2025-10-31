package grpcauth

import (
	"context"
	"errors"
	"strings"

	"github.com/syntaxfa/quick-connect/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type TokenValidator interface {
	ValidateToken(tokenString string) (*types.UserClaims, error)
}

// RoleManager manager method access.
type RoleManager interface {
	GetRequireRoles(method string) []types.Role
}

func NewAuthInterceptor(validator TokenValidator, manager RoleManager) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requiredRoles := manager.GetRequireRoles(info.FullMethod)

		if len(requiredRoles) == 0 {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		values := md.Get(types.AuthorizationKey)
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		tokenString := strings.TrimPrefix(values[0], "Bearer ")

		claims, vErr := validator.ValidateToken(tokenString)
		if vErr != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", vErr)
		}

		hasPermission := false
		for _, role := range claims.Roles {
			for _, requireRole := range requiredRoles {
				if role == requireRole {
					hasPermission = true

					break
				}
			}
		}
		if !hasPermission {
			return nil, status.Errorf(codes.Unauthenticated, "permission denied")
		}

		newCtx := context.WithValue(ctx, types.UserContextKey, claims)

		return handler(newCtx, req)
	}
}

func ExtractUserClaimsFromContext(ctx context.Context) (*types.UserClaims, error) {
	userClaims, ok := ctx.Value(types.UserContextKey).(*types.UserClaims)
	if !ok {
		return nil, errors.New("user claims not provided")
	}

	return userClaims, nil
}
