package http

import (
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

// chatWSHandler docs
//
//	@Summary		chat websocket
//	@Description	chat websocket
//	@Tags			Chats
//	@Accept			json
//	@Produce		json
//	@Router			/chats/clients [GET].
func (h Handler) chatWSHandler(c echo.Context) error {
	conn, uErr := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if uErr != nil {
		h.logger.Error("failed to upgrade connection",
			slog.String("error", uErr.Error()),
			slog.String("remote_addr", c.Request().RemoteAddr))

		return echo.NewHTTPError(http.StatusInternalServerError, "could not upgrade connection")
	}

	for {
		messageType, data, rErr := conn.ReadMessage()
		if rErr != nil {
			h.logger.Error(rErr.Error())
		}

		if wErr := conn.WriteMessage(messageType, data); wErr != nil {
			h.logger.Error(wErr.Error())
		}
	}
}
