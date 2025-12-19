package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "everything is ok!!!"})
}
