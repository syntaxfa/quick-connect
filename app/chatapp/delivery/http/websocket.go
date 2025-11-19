package http

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/auth"
)

func (h Handler) upgradeToWebsocket(c echo.Context) error {
	const op = "http.upgradeToWebsocket"

	claims, cErr := auth.GetUserClaimFormContext(c)
	if cErr != nil {
		h.logger.WarnContext(c.Request().Context(), "unauthorized websocket attempt", slog.String("op", op),
			slog.String("error", cErr.Error()))

		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	responseHeader := http.Header{}
	responseHeader.Set("Sec-WebSocket-Protocol", c.Request().Header.Get("Sec-WebSocket-Protocol"))

	conn, uErr := h.upgrader.Upgrade(c.Response(), c.Request(), responseHeader)
	if uErr != nil {
		h.logger.ErrorContext(c.Request().Context(), "could not upgrade connection", slog.String("op", op),
			slog.String("error", uErr.Error()))

		return echo.NewHTTPError(http.StatusNotAcceptable, "could not upgrade connection")
	}

	h.svc.HandleNewConnection(h.appCtx, conn, claims.UserID)

	return nil
}
