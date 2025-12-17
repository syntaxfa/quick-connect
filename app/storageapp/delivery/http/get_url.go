package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

// getPublicDownloadLink docs
// @Security JWT
// @Router /files/{fileID} [GET]
// @Summary get public link
// @Description get public link
// @Tags Storage
// @Accept json
// @Produce json
// @Param fileID path string true "file ID"
// @Failure 404 {string} conversation does not exist
// @Failure 500 {string} something went wrong.
func (h Handler) getPublicLink(c echo.Context) error {
	resp, sErr := h.svc.GetPublicLink(c.Request().Context(), types.ID(c.Param("fileID")))
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, resp)
}
