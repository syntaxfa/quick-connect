package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/app/adminapp/service"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"github.com/syntaxfa/quick-connect/protobuf/shared/golang/errdetailspb"
	"google.golang.org/grpc/status"
)

type Handler struct {
	logger          *slog.Logger
	t               *translation.Translate
	authSvc         service.AuthService
	userSvc         service.UserService
	conversationSvc service.ConversationService
	chatWSURL       string
}

func NewHandler(logger *slog.Logger, t *translation.Translate, authSvc service.AuthService, userSvc service.UserService,
	conversationSvc service.ConversationService, chatWSURL string) Handler {
	return Handler{
		t:               t,
		logger:          logger,
		authSvc:         authSvc,
		userSvc:         userSvc,
		conversationSvc: conversationSvc,
		chatWSURL:       chatWSURL,
	}
}

func (h Handler) renderErrorPartial(c echo.Context, httpStatus int, errorContent string) error {
	msg := h.t.TranslateMessage(errorContent)

	html := `<div id="modal-error-message" class="error" hx-swap-oob="true">` + msg + `</div>`
	html += `<div id="error-message" class="error" hx-swap-oob="true">` + msg + `</div>`

	// TODO: HTMX can't handle errors when status code is not 200
	c.Response().Header().Set("X-Http-Status", strconv.Itoa(httpStatus))
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
		if d, assertOk := detail.(*errdetailspb.BadRequest); assertOk {
			var htmlBuilder strings.Builder
			htmlBuilder.WriteString("<ul class='error-list'>")

			for _, v := range d.GetFieldViolations() {
				translatedDesc := h.t.TranslateMessage(v.GetDescription())
				htmlBuilder.WriteString(fmt.Sprintf("<li>%s: %s</li>", v.GetField(), translatedDesc))
			}

			htmlBuilder.WriteString("</ul>")
			finalErrorMessage = htmlBuilder.String()
		}
	}

	return h.renderErrorPartial(c, httpStatus, finalErrorMessage)
}

func (h Handler) ShowSuccessToast(c echo.Context) error {
	message := c.QueryParam("message")
	if message == "" {
		message = "Success!"
	}

	data := map[string]interface{}{
		"Message":   message,
		"Timestamp": time.Now().UnixNano(),
	}
	return c.Render(http.StatusOK, "toast_success", data)
}
