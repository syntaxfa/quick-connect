package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "everything is ok")
}
