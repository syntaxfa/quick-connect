package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/jwtvalidator"
	"github.com/syntaxfa/quick-connect/types"
)

const UserContextKey = "user_claims"

type Middleware struct {
	validator *jwtvalidator.Validator
}

func New(validator *jwtvalidator.Validator) *Middleware {
	return &Middleware{
		validator: validator,
	}
}

// RequireAuth (Authentication)
func (m *Middleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		tokenStr, eErr := extractToken(authHeader)
		if eErr != nil {
			return c.JSON(http.StatusUnauthorized, eErr.Error())
		}

		claims, jErr := m.validator.ValidateToken(tokenStr)
		if jErr != nil {
			return c.JSON(http.StatusUnauthorized, jErr.Error())
		}

		c.Set(UserContextKey, claims)

		return next(c)
	}
}

func (m *Middleware) RequireRole(roles []types.Role) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userClaims, gErr := GetUserClaimFormContext(c)
			if gErr != nil {
				return echo.NewHTTPError(http.StatusForbidden, map[string]string{"error": gErr.Error()})
			}

			for _, userRole := range userClaims.Roles {
				for _, requireRole := range roles {
					if userRole == requireRole {
						return next(c)
					}
				}
			}

			return c.JSON(http.StatusForbidden, "You do not have access")
		}
	}
}

func extractToken(authHeader string) (string, error) {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header")
	}

	return parts[1], nil
}

func GetUserClaimFormContext(c echo.Context) (types.UserClaims, error) {
	claimsStr := c.Get(UserContextKey)
	if claimsStr == "" {
		return types.UserClaims{}, errors.New("user claims not found in context")
	}

	userClaims, ok := claimsStr.(*types.UserClaims)
	if !ok {
		return types.UserClaims{}, errors.New("invalid user claims")
	}

	return *userClaims, nil
}
