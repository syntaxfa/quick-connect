package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetUser docs
//
//	@Summary		get user
//	@Description	get user des
//	@Tags			Micro1
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string "everything is good"
//	@Failure		500	{string}	something	went	wrong
//	@Router			/user [GET].
func (h Handler) GetUser(c echo.Context) error {
	if sErr := h.svc.GetUser(c.Request().Context()); sErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, sErr.Error())
	}

	return c.JSON(http.StatusOK, "ok")
}
