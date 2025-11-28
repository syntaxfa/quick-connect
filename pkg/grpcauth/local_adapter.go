package grpcauth

import (
	"context"
	"strings"

	"github.com/syntaxfa/quick-connect/pkg/jwtvalidator"
	"github.com/syntaxfa/quick-connect/pkg/rolemanager"
	"github.com/syntaxfa/quick-connect/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Protect acts as a security gatekeeper for local adapters, simulating the behavior of a gRPC Interceptor.
// It validates the authorization token extracted from the context, checks role permissions against the RoleManager,
// and returns the UserClaims if the request is authorized.
//
// This function is crucial in the "Code-level Monolith" architecture where network layers (and thus gRPC Interceptors)
// are bypassed. It ensures that direct function calls between services remain secure.
//
// Parameters:
//   - ctx: The context containing the authorization token (expected in types.AuthorizationKey).
//   - roleManager: The component defining which roles are required for a given method.
//   - jwtValidator: The component responsible for validating the JWT token string.
//   - fullMethod: The gRPC method name (e.g., "/manager.UserService/UserProfile") used to look up required roles.
//
// Returns:
//   - *types.UserClaims: The extracted user claims if authentication and authorization succeed.
//     Returns nil if the method is public (requires no roles).
//   - error: An error if the token is missing, invalid, or the user lacks permission.
func Protect(ctx context.Context, roleManager *rolemanager.RoleManager, jwtValidator *jwtvalidator.Validator,
	fullMethod string) (*types.UserClaims, error) {
	requiredRoles := roleManager.GetRequireRoles(fullMethod)

	rawToken, ok := ctx.Value(types.AuthorizationKey).(string)
	if !ok || rawToken == "" {
		return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	tokenString := strings.TrimPrefix(rawToken, "Bearer ")

	claims, vErr := jwtValidator.ValidateToken(tokenString)
	if vErr != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %s", vErr.Error())
	}

	hasPermission := false
	for _, role := range claims.Roles {
		for _, reqRole := range requiredRoles {
			if role == reqRole {
				hasPermission = true

				break
			}
		}
	}
	if !hasPermission {
		return nil, status.Errorf(codes.Unauthenticated, "permission denied")
	}

	return claims, nil
}
