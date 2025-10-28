package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
)

func (h Handler) ShowLoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login_layout", nil)
}

func (h Handler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	loginReq := &authpb.LoginRequest{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
	}

	if loginReq.Username == "" || loginReq.Password == "" {
		return h.renderErrorPartial(c, http.StatusBadRequest, h.t.TranslateMessage(servermsg.MsgUsernameAndPasswordAreRequired))
	}

	loginResp, err := h.authAd.Login(ctx, loginReq)

	if err != nil {
		return h.renderGRPCError(c, "gRPC login call failed", err)
	}

	accessMaxAge := int(loginResp.GetAccessExpiresIn())
	refreshMaxAge := int(loginResp.GetRefreshExpiresIn())

	accessCookie := new(http.Cookie)
	accessCookie.Name = "access_token"
	accessCookie.Value = loginResp.GetAccessToken()
	accessCookie.Path = "/"
	accessCookie.HttpOnly = true
	accessCookie.Secure = c.Scheme() == "https"
	accessCookie.SameSite = http.SameSiteLaxMode
	accessCookie.MaxAge = accessMaxAge
	c.SetCookie(accessCookie)

	refreshCookie := new(http.Cookie)
	refreshCookie.Name = "refresh_token"
	refreshCookie.Value = loginResp.GetRefreshToken()
	refreshCookie.Path = "/"
	refreshCookie.HttpOnly = true
	refreshCookie.Secure = c.Scheme() == "https"
	refreshCookie.SameSite = http.SameSiteLaxMode
	refreshCookie.MaxAge = refreshMaxAge
	c.SetCookie(refreshCookie)

	c.Response().Header().Set("HX-Redirect", "/dashboard")
	return c.NoContent(http.StatusOK)
}
