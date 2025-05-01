package servermsg

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"net/http"
)

type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func mapKindToHTTPStatusCode(kind richerror.Kind) int {
	switch kind {
	case richerror.KindInvalid:
		return http.StatusUnprocessableEntity
	case richerror.KindUnAuthorized:
		return http.StatusUnauthorized
	case richerror.KindNotFound:
		return http.StatusNotFound
	case richerror.KindForbidden:
		return http.StatusForbidden
	case richerror.KindBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func HTTPMsg(c echo.Context, err error, t *translation.Translate) error {
	var serverErrCode = http.StatusInternalServerError

	var message string
	var code int
	var errFields map[string]string

	var richErr richerror.RichError
	if errors.As(err, &richErr) {
		message = t.TranslateMessage(richErr.Message())

		code = mapKindToHTTPStatusCode(richErr.Kind())
		if code >= serverErrCode {
			message = MsgSomethingWentWrong
		}

		errFields = richErr.ErrorFields()
	} else {
		message, code = MsgSomethingWentWrong, serverErrCode
	}

	data := ErrorResponse{
		Message: message,
		Errors:  errFields,
	}

	return c.JSON(code, data)
}
