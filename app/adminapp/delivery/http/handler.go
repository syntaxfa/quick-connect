package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/adapter/manager"
	"github.com/syntaxfa/quick-connect/pkg/translation"
)

type Handler struct {
	logger *slog.Logger
	t      *translation.Translate
	authAd *manager.AuthAdapter
}

func NewHandler(logger *slog.Logger, t *translation.Translate, authAd *manager.AuthAdapter) Handler {
	return Handler{
		t:      t,
		logger: logger,
		authAd: authAd,
	}
}

func (h Handler) renderErrorPartial(c echo.Context, httpStatus int, errorMessage string) error {
	html := `<div id="error-message" class="error">` + h.t.TranslateMessage(errorMessage) + `</div>`

	// TODO: HTMX can't handle errors when status code is not 200
	c.Response().Header().Set("X-HTTP-Status", fmt.Sprintf("%d", httpStatus))
	//return c.HTML(httpStatus, html)
	return c.HTML(http.StatusOK, html)
}

func (h Handler) logError(c echo.Context, err error, message string) {
	h.logger.Error(
		message,
		slog.String("error", err.Error()),
		slog.String("method", c.Request().Method),
		slog.String("path", c.Path()),
	)
}
