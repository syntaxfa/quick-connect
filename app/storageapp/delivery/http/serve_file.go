package http

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h Handler) ServeFile(c echo.Context) error {
	relativePath := c.Param("*")

	if relativePath == "" || strings.Contains(relativePath, "..") {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid path"})
	}

	fullPath := filepath.Join(h.localRootPath, relativePath)

	return c.File(fullPath)
}
