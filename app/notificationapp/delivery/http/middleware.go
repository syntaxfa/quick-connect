package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

type Response struct {
	UserID string `json:"user_id"`
}

func validateExternalToken(getUserIDURL string, logger *slog.Logger, httpClient *http.Client) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			const op = "delivery.middleware.validateExternalToken"

			token := c.Request().Header.Get("Identify-Token")
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "identify token is required")
			}

			data := map[string]string{
				"token": token,
			}

			jsonData, mErr := json.Marshal(data)
			if mErr != nil {
				errlog.WithoutErr(richerror.New(op).WithWrapError(mErr).WithKind(richerror.KindUnexpected), logger)

				return echo.NewHTTPError(http.StatusInternalServerError)
			}

			// #nosec G107
			req, reqErr := http.NewRequestWithContext(c.Request().Context(), http.MethodPost,
				getUserIDURL, bytes.NewBuffer(jsonData))
			if reqErr != nil {
				errlog.WithoutErr(richerror.New(op).WithWrapError(reqErr).WithKind(richerror.KindUnexpected), logger)

				return echo.NewHTTPError(http.StatusInternalServerError)
			}

			req.Header.Set("Content-Type", "application/json")

			resp, pErr := httpClient.Do(req)
			if pErr != nil {
				errlog.WithoutErr(richerror.New(op).WithWrapError(pErr).WithKind(richerror.KindUnexpected), logger)

				return echo.NewHTTPError(http.StatusInternalServerError)
			}

			defer func() {
				if cErr := resp.Body.Close(); cErr != nil {
					errlog.WithoutErr(richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), logger)
				}
			}()

			body, rErr := io.ReadAll(resp.Body)
			if rErr != nil {
				errlog.WithoutErr(richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected), logger)

				return echo.NewHTTPError(http.StatusInternalServerError)
			}

			if resp.Status != "200 OK" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Identify token is not valid")
			}

			fmt.Println(string(body))

			var response Response
			uErr := json.Unmarshal(body, &response)
			if uErr != nil {
				errlog.WithoutErr(richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected), logger)

				return echo.NewHTTPError(http.StatusInternalServerError)
			}

			c.Set("user_id", response.UserID)

			return next(c)
		}
	}
}
