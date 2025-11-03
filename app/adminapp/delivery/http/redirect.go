package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func redirectToDashboard(c echo.Context) error {
	//nolint:goconst // It's ok
	isHTMX := c.Request().Header.Get("Hx-Request") == "true"

	if isHTMX {
		c.Response().Header().Set("Hx-Redirect", "/dashboard")

		return c.NoContent(http.StatusOK)
	}

	return c.Redirect(http.StatusSeeOther, "/dashboard")
}

func redirectToLogin(c echo.Context) error {
	isHTMX := c.Request().Header.Get("Hx-Request") == "true"

	if isHTMX {
		c.Response().Header().Set("Hx-Redirect", "/login")

		return c.NoContent(http.StatusOK)
	}

	return c.Redirect(http.StatusSeeOther, "/login")
}
