package http

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/app/storageapp/service"
	"github.com/syntaxfa/quick-connect/pkg/auth"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

// UploadFile docs
// @Summary Upload a file
// @Description Upload a file to storage (S3 or Local)
// @Tags Storage
// @Accept mpfd
// @Produce json
// @Param file formData file true "The file to upload"
// @Param is_public formData boolean false "Is the file public?"
// @Success 201 {object} service.File
// @Failure 400 {string} string "File is required"
// @Failure 401 {string} string "Unauthorized"
// @Failure 413 {string} string "File size limit exceeded"
// @Failure 500 {string} string something went wrong
// @Security JWT
// @Router /files [POST].
func (h Handler) upload(c echo.Context) error {
	claims, cErr := auth.GetUserClaimFormContext(c)
	if cErr != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, cErr.Error())
	}

	fileHeader, formErr := c.FormFile("file")
	if formErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "file is required")
	}

	if fileHeader.Size > h.maxSize {
		return echo.NewHTTPError(http.StatusRequestEntityTooLarge, "file size limit exceeded")
	}

	src, oErr := fileHeader.Open()
	if oErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open file")
	}
	defer func() {
		if closeErr := src.Close(); closeErr != nil {
			h.logger.Error("can't close file source", slog.String("error", closeErr.Error()))
		}
	}()

	isPublic, pErr := strconv.ParseBool(c.FormValue("is_public"))
	if pErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "is_public is not valid")
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	resp, sErr := h.svc.Upload(c.Request().Context(), service.UploadRequest{
		UploaderID:  claims.UserID,
		File:        src,
		Filename:    fileHeader.Filename,
		Size:        fileHeader.Size,
		ContentType: contentType,
		IsPublic:    isPublic,
	})

	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusCreated, resp)
}
