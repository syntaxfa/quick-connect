package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func redirectToDashboard(c echo.Context) error {
	if isHTMX(c) {
		setHxRedirect(c, "/dashboard")

		return c.NoContent(http.StatusOK)
	}

	return c.Redirect(http.StatusSeeOther, "/dashboard")
}

func redirectToLogin(c echo.Context) error {
	if isHTMX(c) {
		setHxRedirect(c, "/login")

		return c.NoContent(http.StatusOK)
	}

	return c.Redirect(http.StatusSeeOther, "/login")
}
