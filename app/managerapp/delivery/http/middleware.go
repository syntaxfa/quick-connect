package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const APIKey = "APIKey"

func AccessToPrivateEndpoint(apiKey string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqAPIKey := c.QueryParam(APIKey)

			if reqAPIKey == "" {
				return echo.NewHTTPError(http.StatusForbidden, "API Key not provided")
			}

			if reqAPIKey != apiKey {
				return echo.NewHTTPError(http.StatusForbidden, "API Key is not valid")
			}

			return next(c)
		}
	}
}
