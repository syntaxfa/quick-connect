package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/app/storyapp/service"
	"github.com/syntaxfa/quick-connect/pkg/auth"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

// AddStory docs
// @Summary Add Story
// @Description adding new story
// @Tags Story
// @Accept json
// @Produce json
// @Param Request body service.AddStoryRequest true "create and return story"
// @Success 201 {object} service.AddStoryResponse
// @Failure 401 {string} unauthorized
// @Failure 422 {object} servermsg.ErrorResponse
// @Failure 500 {string} something went wrong
// @Security JWT
// @Router /stories [POST].
func (h Handler) AddStory(c echo.Context) error {
	var req service.AddStoryRequest
	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	claims, cErr := auth.GetUserClaimFormContext(c)
	if cErr != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, cErr.Error())
	}

	resp, sErr := h.svc.AddStory(c.Request().Context(), req, claims.UserID)
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusCreated, resp)
}
