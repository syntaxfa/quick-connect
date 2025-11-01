package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func redirectToDashboard(c echo.Context) error {
	isHTMX := c.Request().Header.Get("HX-Request") == "true"

	if isHTMX {
		c.Response().Header().Set("HX-Redirect", "/dashboard")

		return c.NoContent(http.StatusOK)
	} else {
		return c.Redirect(http.StatusSeeOther, "/dashboard")
	}
}

func redirectToLogin(c echo.Context) error {
	isHTMX := c.Request().Header.Get("HX-Request") == "true"

	if isHTMX {
		c.Response().Header().Set("HX-Redirect", "/login")

		return c.NoContent(http.StatusOK)
	} else {
		return c.Redirect(http.StatusSeeOther, "/login")
	}
}
