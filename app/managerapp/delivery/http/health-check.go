package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// healthCheck docs
//
//	@Summary		health check
//	@Description	health check managers service
//	@Tags			manager
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string "everything is good"
//	@Failure		500	{string}	something	went	wrong
//	@Router			/health-check [GET].
func (h Handler) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "everything is ok!!!")
}
