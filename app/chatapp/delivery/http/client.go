package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// WSClientHandler docs
//
//	@Summary		client chat websocket
//	@Description	client chat websocket
//	@Tags			Websocket
//	@Accept			json
//	@Produce		json
//	@Router			/chats/clients [GET].
func (h Handler) WSClientHandler(c echo.Context) error {
	conn, uErr := h.upgrader.Upgrade(c.Response(), c.Request())
	if uErr != nil {
		return echo.NewHTTPError(http.StatusNotAcceptable, "could not upgrade connection")
	}

	h.svc.JoinClient(conn, "guest")

	return nil
}
