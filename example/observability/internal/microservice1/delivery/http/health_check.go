package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// healthCheck docs
//
//	@Summary		health check chat service
//	@Description	health check chat service
//	@Tags			Micro1
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string "everything is good"
//	@Failure		500	{string}	something	went	wrong
//	@Router			/health-check [GET].
func (h Handler) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "everything is ok!!!"})
}
