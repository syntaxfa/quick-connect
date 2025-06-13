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

	userID, ok := c.Get("user_id").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user id is not valid")
	}

	h.svc.JoinClient(c.Request().Context(), conn, userID)

	return nil
}
