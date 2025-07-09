package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

// updateTemplate docs
// @Router /v1/templates/{templateID} [PUT]
// @Summary update template
// @Description This API endpoint updates a specific template.
// @Tags NotificationAdmin
// @Accept json
// @Produce json
// @Param Request body service.AddTemplateRequest true "template"
// @Param templateID path string true "ID of the template to update"
// @Success 200 {object} service.Template
// @Failure 400 {string} string Bad Request
// @Failure 404 {string} the template with this templateID does not exist
// @Failure 409 {string} the name of template has conflict
// @Failure 422 {object} servermsg.ErrorResponse
// @Failure 500 {string} something went wrong.
func (h Handler) updateTemplate(c echo.Context) error {
	var req service.AddTemplateRequest
	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	resp, sErr := h.svc.UpdateTemplate(c.Request().Context(), types.ID(c.Param("templateID")), req)
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, resp)
}
