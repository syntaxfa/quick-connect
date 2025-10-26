package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) ShowLoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login_layout", nil)
}
