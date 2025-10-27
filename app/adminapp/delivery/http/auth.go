package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return h.renderErrorPartial(c, http.StatusBadRequest, "Username and password are required")
	}

	loginResp, err := h.authAd.Login(ctx, loginReq)

	if err != nil {
		h.logError(c, err, "gRPC login call failed")

		st, ok := status.FromError(err)
		if !ok {
			return h.renderErrorPartial(c, http.StatusInternalServerError, "An unexpected error occurred")
		}

		errorMessage := st.Message()
		httpStatus := http.StatusInternalServerError // پیش‌فرض

		switch st.Code() {
		case codes.InvalidArgument:
			httpStatus = http.StatusBadRequest
		case codes.Unauthenticated:
			httpStatus = http.StatusUnauthorized
		case codes.NotFound:
			httpStatus = http.StatusUnauthorized // for security is better return 401 for 404
		}

		return h.renderErrorPartial(c, httpStatus, errorMessage)
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

	c.Response().Header().Set("HX-Redirect", "/dashboard") // ریدایرکت به داشبورد
	return c.NoContent(http.StatusOK)
}
