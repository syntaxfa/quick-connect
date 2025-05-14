package errlog

import (
	"errors"
	"log/slog"

	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

func ErrLog(err error, logger *slog.Logger) error {
	var richErr richerror.RichError

	if errors.As(err, &richErr) {
		logger.Error(richErr.Message(), slog.String("operation", richErr.Operation()),
			slog.Int("kind", int(richErr.Kind())), slog.Any("detail", richErr.ExtraDetail()))
	}

	logger.Error(err.Error())

	return err
}

func WithoutErr(err error, logger *slog.Logger) {
	var richErr richerror.RichError

	if errors.As(err, &richErr) {
		logger.Error(richErr.Message(), slog.String("operation", richErr.Operation()),
			slog.Int("kind", int(richErr.Kind())), slog.Any("detail", richErr.ExtraDetail()))
	}

	logger.Error(err.Error())
}
