package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const (
	SecureSchema = "https"
)

func isUserHaveAuthCookie(c echo.Context, logger *slog.Logger) bool {
	cookie, err := c.Cookie(string(types.TokenTypeRefresh))

	if err != nil {
		if !errors.Is(err, http.ErrNoCookie) {
			errlog.WithoutErr(richerror.New("isUserHaveAuthCookie").WithKind(richerror.KindUnexpected).WithWrapError(err), logger)
		}

		clearAuthCookie(c)

		return false
	}

	if cookie.Value == "" {
		clearAuthCookie(c)

		return false
	}

	return true
}

func setAuthCookie(c echo.Context, accessToken, refreshToken string, accessExpiry, refreshExpiry int) {
	accessCookie := new(http.Cookie)
	accessCookie.Name = string(types.TokenTypeAccess)
	accessCookie.Value = accessToken
	accessCookie.Path = "/"
	accessCookie.HttpOnly = true
	accessCookie.Secure = c.Scheme() == SecureSchema
	accessCookie.SameSite = http.SameSiteLaxMode
	accessCookie.MaxAge = accessExpiry
	c.SetCookie(accessCookie)

	refreshCookie := new(http.Cookie)
	refreshCookie.Name = string(types.TokenTypeRefresh)
	refreshCookie.Value = refreshToken
	refreshCookie.Path = "/"
	refreshCookie.HttpOnly = true
	refreshCookie.Secure = c.Scheme() == SecureSchema
	refreshCookie.SameSite = http.SameSiteLaxMode
	refreshCookie.MaxAge = refreshExpiry
	c.SetCookie(refreshCookie)
}

func clearAuthCookie(c echo.Context) {
	accessCookie := new(http.Cookie)
	accessCookie.Name = string(types.TokenTypeAccess)
	accessCookie.Value = ""
	accessCookie.Path = "/"
	accessCookie.HttpOnly = true
	accessCookie.Secure = c.Scheme() == SecureSchema
	accessCookie.SameSite = http.SameSiteLaxMode
	accessCookie.MaxAge = -1
	c.SetCookie(accessCookie)

	refreshCookie := new(http.Cookie)
	refreshCookie.Name = string(types.TokenTypeRefresh)
	refreshCookie.Value = ""
	refreshCookie.Path = "/"
	refreshCookie.HttpOnly = true
	refreshCookie.Secure = c.Scheme() == SecureSchema
	refreshCookie.SameSite = http.SameSiteLaxMode
	refreshCookie.MaxAge = -1
	c.SetCookie(refreshCookie)
}

func getAccessTokenFromCookie(c echo.Context, logger *slog.Logger) (string, bool) {
	cookie, err := c.Cookie(string(types.TokenTypeAccess))

	if err != nil {
		if !errors.Is(err, http.ErrNoCookie) {
			errlog.WithoutErr(richerror.New("isUserHaveAuthCookie").WithKind(richerror.KindUnexpected).WithWrapError(err), logger)
		}

		return "", false
	}

	if cookie.Value == "" {
		return "", false
	}

	return cookie.Value, true
}

func getRefreshTokenFromCookie(c echo.Context, logger *slog.Logger) (string, bool) {
	cookie, err := c.Cookie(string(types.TokenTypeRefresh))

	if err != nil {
		if !errors.Is(err, http.ErrNoCookie) {
			errlog.WithoutErr(richerror.New("isUserHaveAuthCookie").WithKind(richerror.KindUnexpected).WithWrapError(err), logger)
		}

		return "", false
	}

	if cookie.Value == "" {
		return "", false
	}

	return cookie.Value, true
}
