package http

import (
	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/auth"

	"log/slog"
	"net/http"
)

func (h Handler) upgradeToWebsocket(c echo.Context) error {
	const op = "http.upgradeToWebsocket"

	claims, cErr := auth.GetUserClaimFormContext(c)
	if cErr != nil {
		h.logger.WarnContext(c.Request().Context(), "unauthorized websocket attempt", slog.String("op", op),
			slog.String("error", cErr.Error()))

		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	conn, uErr := h.upgrader.Upgrade(c.Response(), c.Request())
	if uErr != nil {
		h.logger.ErrorContext(c.Request().Context(), "could not upgrade connection", slog.String("op", op),
			slog.String("error", uErr.Error()))

		return echo.NewHTTPError(http.StatusNotAcceptable, "could not upgrade connection")
	}

	h.svc.HandleNewConnection(h.appCtx, conn, claims.UserID)

	return nil
}
