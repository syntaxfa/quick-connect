package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) wsNotification(c echo.Context) error {
	conn, uErr := h.upgrader.Upgrade(c.Response(), c.Request())
	if uErr != nil {
		return echo.NewHTTPError(http.StatusNotAcceptable, "could not upgrade connection")
	}

	h.svc.JoinClient(c.Request().Context(), conn, "1")

	return nil
}
