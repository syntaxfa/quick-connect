package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

// createTemplate docs
// @Summary create new template
// @Description This API endpoint creates a new template.
// @Tags NotificationAdmin
// @Accept json
// @Produce json
// @Param Request body service.AddTemplateRequest true "template"
// @Success 201 {object} service.Template
// @Failure 400 {string} string Bad Request
// @Failure 409 {string} the name of template has conflict
// @Failure 422 {object} servermsg.ErrorResponse
// @Failure 500 {string} something went wrong
// @Router /v1/templates [POST].
func (h Handler) createTemplate(c echo.Context) error {
	var req service.AddTemplateRequest
	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	resp, sErr := h.svc.AddTemplate(c.Request().Context(), req)
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusCreated, resp)
}
