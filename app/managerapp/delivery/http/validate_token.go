package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

// ValidateToken docs
// @Summary ValidateToken JWT
// @Description This API endpoint validates a JSON Web Token (JWT) to ensure its authenticity and integrity. It checks the token's signature, expiration, and claims.
// @Tags Token
// @Accept json
// @Produce json
// @Param Request body tokenservice.TokenVerifyRequest true "check token validation"
// @Success 200 "jwt token is valid"
// @Failure 400 {string} string Bad Request
// @Failure 401 {string} string "invalid or expired jwt"
// @Failure 422 {object} servermsg.ErrorResponse
// @Failure 500 {string} something went wrong
// @Router /tokens/validate [POST].
func (h Handler) ValidateToken(c echo.Context) error {
	var req tokenservice.TokenVerifyRequest

	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	resp, sErr := h.tokenSvc.ValidateToken(req.Token)
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, resp)
}
