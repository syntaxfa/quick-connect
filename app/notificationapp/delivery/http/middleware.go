package http

import (
	"bytes"
	"encoding/json"
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

			userID, err := fetchUserIDFromToken(c, token, getUserIDURL, httpClient, logger, op)
			if err != nil {
				return err
			}

			c.Set("user_id", userID)
			return next(c)
		}
	}
}

func fetchUserIDFromToken(c echo.Context, token, getUserIDURL string, httpClient *http.Client, logger *slog.Logger,
	op string) (string, error) {
	jsonData, err := createTokenRequest(token, logger, op)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError)
	}

	//nolint:bodyclose // Body is closed correctly in the 'closeResponseBody' helper function
	resp, err := sendTokenValidationRequest(c, getUserIDURL, jsonData, httpClient, logger, op)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError)
	}

	defer closeResponseBody(resp.Body, logger, op)

	return parseUserIDFromResponse(resp, logger, op)
}

func createTokenRequest(token string, logger *slog.Logger, op string) ([]byte, error) {
	data := map[string]string{"token": token}

	jsonData, mErr := json.Marshal(data)
	if mErr != nil {
		errlog.WithoutErr(richerror.New(op).WithWrapError(mErr).WithKind(richerror.KindUnexpected), logger)
		return nil, mErr
	}

	return jsonData, nil
}

func sendTokenValidationRequest(c echo.Context, getUserIDURL string, jsonData []byte, httpClient *http.Client,
	logger *slog.Logger, op string) (*http.Response, error) {
	// #nosec G107
	req, reqErr := http.NewRequestWithContext(c.Request().Context(), http.MethodPost, getUserIDURL, bytes.NewBuffer(jsonData))
	if reqErr != nil {
		errlog.WithoutErr(richerror.New(op).WithWrapError(reqErr).WithKind(richerror.KindUnexpected), logger)
		return nil, reqErr
	}

	req.Header.Set("Content-Type", "application/json")

	resp, pErr := httpClient.Do(req)
	if pErr != nil {
		errlog.WithoutErr(richerror.New(op).WithWrapError(pErr).WithKind(richerror.KindUnexpected), logger)
		return nil, pErr
	}

	return resp, nil
}

func closeResponseBody(body io.ReadCloser, logger *slog.Logger, op string) {
	if cErr := body.Close(); cErr != nil {
		errlog.WithoutErr(richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), logger)
	}
}

func parseUserIDFromResponse(resp *http.Response, logger *slog.Logger, op string) (string, error) {
	body, rErr := io.ReadAll(resp.Body)
	if rErr != nil {
		errlog.WithoutErr(richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected), logger)
		return "", echo.NewHTTPError(http.StatusInternalServerError)
	}

	if resp.Status != "200 OK" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Identify token is not valid")
	}

	var response Response
	if uErr := json.Unmarshal(body, &response); uErr != nil {
		errlog.WithoutErr(richerror.New(op).WithWrapError(uErr).WithKind(richerror.KindUnexpected), logger)
		return "", echo.NewHTTPError(http.StatusInternalServerError)
	}

	return response.UserID, nil
}
