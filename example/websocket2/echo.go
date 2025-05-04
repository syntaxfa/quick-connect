package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"io"
)

func main() {
	httpServer()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func httpServer() {
	e := echo.New()

	e.GET("/echo1/", EchoHandler1)
	e.GET("/echo2/", echoHandler2)

	e.Start(":2525")
}

func EchoHandler1(c echo.Context) error {
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

func echoHandler2(c echo.Context) error {
	conn, uErr := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if uErr != nil {
		c.Logger().Error(uErr)
	}

	for {
		messageType, r, err := conn.NextReader()
		if err != nil {
			c.Logger().Error(err)
		}

		w, err := conn.NextWriter(messageType)
		if err != nil {
			c.Logger().Error(err)
		}

		if _, err := io.Copy(w, r); err != nil {
			c.Logger().Error(err)
		}

		if err := w.Close(); err != nil {
			return err
		}
	}
}
