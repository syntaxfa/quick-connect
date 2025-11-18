package http

import (
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
	return h.upgradeToWebsocket(c)
}
