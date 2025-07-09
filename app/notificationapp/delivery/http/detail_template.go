package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

// getDetailTemplate docs
// @Router /v1/templates/{templateID} [GET]
// @Summary retrieve detail template
// @Description This API endpoint retrieve template detail.
// @Tags NotificationAdmin
// @Accept json
// @Produce json
// @Param templateID path string true "ID of the template to update"
// @Success 200 {object} service.Template
// @Failure 404 {string} the template with this templateID does not exist
// @Failure 500 {string} something went wrong.
func (h Handler) getDetailTemplate(c echo.Context) error {
	resp, sErr := h.svc.GetTemplate(c.Request().Context(), types.ID(c.Param("templateID")))
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, resp)
}
