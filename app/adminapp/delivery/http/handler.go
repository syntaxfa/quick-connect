package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/adapter/manager"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"github.com/syntaxfa/quick-connect/protobuf/shared/golang/errdetailspb"
	"google.golang.org/grpc/status"
)

type Handler struct {
	logger *slog.Logger
	t      *translation.Translate
	authAd *manager.AuthAdapter
	userAd *manager.UserAdapter
}

func NewHandler(logger *slog.Logger, t *translation.Translate, authAd *manager.AuthAdapter, userAd *manager.UserAdapter) Handler {
	return Handler{
		t:      t,
		logger: logger,
		authAd: authAd,
		userAd: userAd,
	}
}

func (h Handler) renderErrorPartial(c echo.Context, httpStatus int, errorContent string) error {
	html := `<div id="error-message" class="error">` + h.t.TranslateMessage(errorContent) + `</div>`

	// TODO: HTMX can't handle errors when status code is not 200
	c.Response().Header().Set("X-HTTP-Status", fmt.Sprintf("%d", httpStatus))
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

func (h Handler) renderGRPCError(c echo.Context, operationName string, err error) error {
	h.logError(c, err, operationName)

	st, ok := status.FromError(err)
	if !ok {
		translatedMsg := h.t.TranslateMessage(servermsg.MsgSomethingWentWrong)

		return h.renderErrorPartial(c, http.StatusInternalServerError, translatedMsg)
	}

	httpStatus := servermsg.GRPCCodeToHTTPStatusCode(st.Code())
	finalErrorMessage := h.t.TranslateMessage(st.Message())

	for _, detail := range st.Details() {
		switch d := detail.(type) {
		case *errdetailspb.BadRequest:
			var htmlBuilder strings.Builder
			htmlBuilder.WriteString("<ul class='error-list'>")

			for _, v := range d.GetFieldViolations() {
				translatedDesc := h.t.TranslateMessage(v.GetDescription())
				htmlBuilder.WriteString(fmt.Sprintf("<li>%s: %s</li>", v.GetField(), translatedDesc))
			}

			htmlBuilder.WriteString("</ul>")
			finalErrorMessage = htmlBuilder.String()
			break
		}
	}

	return h.renderErrorPartial(c, httpStatus, finalErrorMessage)
}
