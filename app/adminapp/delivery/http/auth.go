package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
)

func (h Handler) ShowLoginPage(c echo.Context) error {
	if exist := isUserHaveAuthCookie(c, h.logger); exist {
		return redirectToDashboard(c)
	}

	return c.Render(http.StatusOK, "login_layout", nil)
}

func (h Handler) Login(c echo.Context) error {
	if exist := isUserHaveAuthCookie(c, h.logger); exist {
		return redirectToDashboard(c)
	}

	ctx := c.Request().Context()

	loginReq := &authpb.LoginRequest{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
	}

	if loginReq.GetUsername() == "" || loginReq.GetPassword() == "" {
		return h.renderErrorPartial(c, http.StatusBadRequest, h.t.TranslateMessage(servermsg.MsgUsernameAndPasswordAreRequired))
	}

	loginResp, err := h.authAd.Login(ctx, loginReq)

	if err != nil {
		return h.renderGRPCError(c, "gRPC login call failed", err)
	}

	setAuthCookie(c, loginResp.GetAccessToken(), loginResp.GetRefreshToken(), int(loginResp.GetAccessExpiresIn()),
		int(loginResp.GetRefreshExpiresIn()))

	return redirectToDashboard(c)
}

// ShowLogoutConfirm renders the logout confirmation modal.
func (h Handler) ShowLogoutConfirm(c echo.Context) error {
	return c.Render(http.StatusOK, "logout_confirm_modal", nil)
}

func (h Handler) Logout(c echo.Context) error {
	clearAuthCookie(c)

	c.Response().Header().Set("Hx-Redirect", "/login")

	return c.NoContent(http.StatusOK)
}
