package servermsg

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"github.com/syntaxfa/quick-connect/protobuf/shared/golang/errdetailspb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	case richerror.KindConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func mapKindToGRPCCode(kind richerror.Kind) codes.Code {
	switch kind {
	case richerror.KindInvalid:
		return codes.InvalidArgument
	case richerror.KindUnAuthorized:
		return codes.Unauthenticated
	case richerror.KindNotFound:
		return codes.NotFound
	case richerror.KindForbidden:
		return codes.PermissionDenied
	case richerror.KindBadRequest:
		return codes.InvalidArgument
	case richerror.KindConflict:
		return codes.AlreadyExists
	case richerror.KindUnexpected:
		return codes.Internal
	default:
		return codes.Unknown
	}
}

func HTTPMsg(c echo.Context, err error, t *translation.Translate) error {
	serverErrCode := http.StatusInternalServerError

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

// GRPCMsg converts a richerror into a gRPC status error, potentially with details.
func GRPCMsg(err error, t *translation.Translate, logger *slog.Logger) error {
	unknownCode := codes.Unknown
	internalCode := codes.Internal

	var message string
	var code codes.Code
	var errFields map[string]string
	var richErr richerror.RichError

	if errors.As(err, &richErr) {
		message = t.TranslateMessage(richErr.Message())
		code = mapKindToGRPCCode(richErr.Kind())
		errFields = richErr.ErrorFields()

		logger.Warn("gRPC request failed", "code", code.String(), "operation", richErr.Operation(),
			"detail", richErr.ExtraDetail())

		translationMessage := t.TranslateMessage(message)

		if code == internalCode || code == unknownCode {
			translationMessage = MsgSomethingWentWrong
		}

		if len(errFields) > 0 && code == codes.InvalidArgument {
			badRequestDetails := &errdetailspb.BadRequest{}
			for field, desc := range errFields {
				badRequestDetails.FieldViolations = append(badRequestDetails.FieldViolations,
					&errdetailspb.FieldViolation{
						Field:       field,
						Description: desc,
					})
			}

			st := status.New(code, translationMessage)
			stWithDetails, detailErr := st.WithDetails(badRequestDetails)
			if detailErr != nil {
				logger.Error("failed to add details to gRPC status", slog.String("error", detailErr.Error()))
				return st.Err()
			}

			return stWithDetails.Err()
		}

		return status.Error(code, translationMessage)
	}

	logger.Error("gRPC request failed with unexpected error", slog.String("error", err.Error()))

	return status.Error(internalCode, MsgSomethingWentWrong)
}

func GRPCCodeToHTTPStatusCode(code codes.Code) int {
	switch code {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.NotFound:
		return http.StatusNotFound
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.AlreadyExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
