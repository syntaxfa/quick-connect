package http

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

func isUserHaveAuthCookie(c echo.Context, logger *slog.Logger) bool {
	if cookie, cErr := c.Cookie(string(types.TokenTypeRefresh)); cErr != nil {
		errlog.WithoutErr(richerror.New("Login").WithKind(richerror.KindUnexpected).WithWrapError(cErr), logger)
	} else if vErr := cookie.Valid(); vErr == nil {
		return true
	}

	clearAuthCookie(c)

	return false
}

func setAuthCookie(c echo.Context, accessToken, refreshToken string, accessExpiry, refreshExpiry int) {
	accessCookie := new(http.Cookie)
	accessCookie.Name = string(types.TokenTypeAccess)
	accessCookie.Value = accessToken
	accessCookie.Path = "/"
	accessCookie.HttpOnly = true
	accessCookie.Secure = c.Scheme() == "https"
	accessCookie.SameSite = http.SameSiteLaxMode
	accessCookie.MaxAge = accessExpiry
	c.SetCookie(accessCookie)

	refreshCookie := new(http.Cookie)
	refreshCookie.Name = string(types.TokenTypeRefresh)
	refreshCookie.Value = refreshToken
	refreshCookie.Path = "/"
	refreshCookie.HttpOnly = true
	refreshCookie.Secure = c.Scheme() == "https"
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
	accessCookie.Secure = c.Scheme() == "https"
	accessCookie.SameSite = http.SameSiteLaxMode
	accessCookie.MaxAge = -1
	c.SetCookie(accessCookie)

	refreshCookie := new(http.Cookie)
	refreshCookie.Name = string(types.TokenTypeRefresh)
	refreshCookie.Value = ""
	refreshCookie.Path = "/"
	refreshCookie.HttpOnly = true
	refreshCookie.Secure = c.Scheme() == "https"
	refreshCookie.SameSite = http.SameSiteLaxMode
	refreshCookie.MaxAge = -1
	c.SetCookie(refreshCookie)
}
