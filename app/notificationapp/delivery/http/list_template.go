package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

// ListTemplate docs
// @Router /v1/templates/list [POST]
// @Summary list of all templates
// @Description This API endpoint retrieve all templates.
// @Tags NotificationAdmin
// @Accept json
// @Produce json
// @Param Request body service.ListTemplateRequest true "template list"
// @Success 200 {object} service.ListTemplateResponse
// @Failure 500 {string} something went wrong.
func (h Handler) ListTemplate(c echo.Context) error {
	var req service.ListTemplateRequest
	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	resp, sErr := h.svc.TemplateList(c.Request().Context(), req)
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, resp)
}
