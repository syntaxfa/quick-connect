package ratelimit

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/cachemanager"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

func ByIPAddressMiddleware(cache *cachemanager.CacheManager, maxHint int,
	duration time.Duration, logger *slog.Logger) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			const op = "pkg.ratelimit.middleware.rateLimitMiddleware"
			ctx := c.Request().Context()
			key := getCacheKey(c.RealIP(), c.Path())

			hint, iErr := cache.Incr(ctx, key)
			if iErr != nil {
				errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(iErr).WithKind(richerror.KindUnexpected), logger)

				return echo.NewHTTPError(http.StatusInternalServerError, "rate limit check failed")
			}

			if hint == 1 {
				if expErr := cache.Expire(ctx, key, duration); expErr != nil {
					errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(expErr).WithKind(richerror.KindUnexpected), logger)

					return echo.NewHTTPError(http.StatusInternalServerError, "late limit expiration failed")
				}
			}

			if hint > int64(maxHint) {
				return echo.NewHTTPError(http.StatusTooManyRequests, "Too Many Requests")
			}

			return next(c)
		}
	}
}

func getCacheKey(ipAddress, url string) string {
	return fmt.Sprintf("limits:%s:%s", ipAddress, url)
}
