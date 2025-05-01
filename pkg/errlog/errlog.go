package errlog

import (
	"errors"
	"log/slog"

	"github.com/syntaxfa/quick-connect/pkg/logger"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

func ErrLog(err error) {
	var richErr richerror.RichError

	if errors.As(err, &richErr) {
		logger.L().Error(richErr.Message(), slog.String("operation", richErr.Operation()),
			slog.Int("kind", int(richErr.Kind())), slog.Any("detail", richErr.ExtraDetail()))
	}

	logger.L().Error(err.Error())
}
