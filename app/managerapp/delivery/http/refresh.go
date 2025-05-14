package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

// RefreshToken docs
// @Summary RefreshToken JWT
// @Description This API endpoint refresh JSON Web Token (JWT).
// @Tags Token
// @Accept json
// @Produce json
// @Param Request body tokenservice.TokenRefreshRequest true "generate pair(refresh & access) tokens"
// @Success 200 {object} tokenservice.TokenGenerateResponse
// @Failure 400 {string} string Bad Request
// @Failure 401 {string} string "invalid or expired jwt"
// @Failure 422 {object} servermsg.ErrorResponse
// @Failure 500 {string} something went wrong
// @Router /tokens/refresh [POST].
func (h Handler) RefreshToken(c echo.Context) error {
	var req tokenservice.TokenRefreshRequest

	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	resp, sErr := h.tokenSvc.RefreshTokens(req.RefreshToken)
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, resp)
}
