package http

import (
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
	return h.upgradeToWebsocket(c)
}
