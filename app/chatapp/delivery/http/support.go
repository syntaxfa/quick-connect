package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// WSSupportHandler docs
//
//	@Summary		client chat websocket
//	@Description	client chat websocket
//	@Tags			Websocket
//	@Accept			json
//	@Produce		json
//	@Router			/chats/supports [GET].
func (h Handler) WSSupportHandler(c echo.Context) error {
	conn, uErr := h.upgrader.Upgrade(c.Response(), c.Request())
	if uErr != nil {
		return echo.NewHTTPError(http.StatusNotAcceptable, "could not upgrade connection")
	}

	h.svc.JoinSupport(conn, "support")

	return nil
}
