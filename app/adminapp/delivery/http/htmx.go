package http

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

func isHTMX(c echo.Context) bool {
	return c.Request().Header.Get("Hx-Request") == "true"
}

func setTriggerAfterSettle(c echo.Context, successMsg string) {
	eventDetail := map[string]string{"message": successMsg}
	eventDetailBytes, _ := json.Marshal(eventDetail)

	c.Response().Header().Set("Hx-Trigger-After-Settle", `{"showSuccessToast": `+string(eventDetailBytes)+`}`)
}
