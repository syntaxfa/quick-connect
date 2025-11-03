package http

import "github.com/labstack/echo/v4"

func isHTMX(c echo.Context) bool {
	return c.Request().Header.Get("Hx-Request") == "true"
}
