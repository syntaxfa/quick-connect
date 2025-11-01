package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
)

func (h Handler) ShowLoginPage(c echo.Context) error {
	if exist := isUserHaveAuthCookie(c, h.logger); exist {
		isHTMX := c.Request().Header.Get("HX-Request") == "true"

		if isHTMX {
			c.Response().Header().Set("HX-Redirect", "/dashboard")

			return c.NoContent(http.StatusOK)
		} else {
			return c.Redirect(http.StatusSeeOther, "/dashboard")
		}
	}

	return c.Render(http.StatusOK, "login_layout", nil)
}

func (h Handler) Login(c echo.Context) error {
	if exist := isUserHaveAuthCookie(c, h.logger); exist {
		isHTMX := c.Request().Header.Get("HX-Request") == "true"

		if isHTMX {
			c.Response().Header().Set("HX-Redirect", "/dashboard")

			return c.NoContent(http.StatusOK)
		} else {
			return c.Redirect(http.StatusSeeOther, "/dashboard")
		}
	}

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

	setAuthCookie(c, loginResp.AccessToken, loginResp.RefreshToken, int(loginResp.GetAccessExpiresIn()), int(loginResp.GetRefreshExpiresIn()))

	c.Response().Header().Set("HX-Redirect", "/dashboard")
	return c.NoContent(http.StatusOK)
}
