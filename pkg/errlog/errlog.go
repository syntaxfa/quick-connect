package errlog

import (
	"log/slog"

	"github.com/syntaxfa/quick-connect/pkg/logger"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

func ErrLog(err richerror.RichError) {
	logger.L().Error(err.Message(), slog.String("operation", err.Operation()),
		slog.Int("kind", int(err.Kind())), slog.Any("detail", err.ExtraDetail()))
}
