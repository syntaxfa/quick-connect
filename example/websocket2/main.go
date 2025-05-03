package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func main() {
	upgradeHttpRequestToWebsocket()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func upgradeHttpRequestToWebsocket() {
	e := echo.New()

	e.GET("/", upgradeHTTPRequestToWSHandler)

	e.Start(":2525")
}

func upgradeHTTPRequestToWSHandler(c echo.Context) error {
	ws, uErr := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if uErr != nil {
		c.Logger().Error(uErr)
	}

	for {
		msgType, msg, rErr := ws.ReadMessage()
		if rErr != nil {
			c.Logger().Error(rErr)
		}

		fmt.Printf("message type: %d\n", msgType)

		if wErr := ws.WriteMessage(1, []byte(msg)); wErr != nil {
			c.Logger().Error(wErr)
		}
	}
}
