package http

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/adapter/manager"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/jwtvalidator"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
	"github.com/syntaxfa/quick-connect/types"
)

func setTokenToRequestContextMiddleware(jwtValidator *jwtvalidator.Validator, authAd *manager.AuthAdapter, loginPath string,
	logger *slog.Logger) func(nex echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			const op = "delivery.middleware.setTokenToRequestContext"

			if c.Path() == loginPath {
				return next(c)
			}

			accessToken, acExist := getAccessTokenFromCookie(c, logger)
			if acExist {
				if claims, err := jwtValidator.ValidateToken(accessToken); err == nil {
					setTokenToContext(c, accessToken)
					setUserToContext(c, claims)

					return next(c)
				}
			}

			refreshToken, rtExist := getRefreshTokenFromCookie(c, logger)
			if !rtExist {
				clearAuthCookie(c)

				return redirectToLogin(c)
			}

			token, tErr := authAd.TokenRefresh(c.Request().Context(), &authpb.TokenRefreshRequest{RefreshToken: refreshToken})
			if tErr != nil {
				errlog.WithoutErr(richerror.New(op).WithWrapError(tErr).WithKind(richerror.KindUnexpected).
					WithMessage("refresh token is not valid"), logger)

				clearAuthCookie(c)

				return redirectToLogin(c)
			}

			if claims, err := jwtValidator.ValidateToken(token.GetAccessToken()); err == nil {
				setUserToContext(c, claims)
			}

			setAuthCookie(c, token.GetAccessToken(), token.GetRefreshToken(), int(token.GetAccessExpiresIn()),
				int(token.GetRefreshExpiresIn()))

			setTokenToContext(c, token.GetAccessToken())

			return next(c)
		}
	}
}

func setTokenToContext(c echo.Context, accessToken string) {
	c.Set(string(types.AuthorizationKey), "Bearer "+accessToken)
}

func setUserToContext(c echo.Context, claims *types.UserClaims) {
	user := convertClaimsToUser(claims)

	c.Set("User", user)
}

func getUserFromContext(c echo.Context) (User, bool) {
	user, ok := c.Get("User").(User)

	return user, ok
}
