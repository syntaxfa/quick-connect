package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"time"
)

/*
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)


   1006

      1006 is a reserved value and MUST NOT be set as a status code in a
      Close control frame by an endpoint.  It is designated for use in
      applications expecting a status code to indicate that the
      connection was closed abnormally, e.g., without sending or
      receiving a Close control frame.
*/

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	e := echo.New()

	e.GET("/messages/", messagesHandler)

	e.Start(":2525")
}

func messagesHandler(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err)
	}

	closeChan := make(chan struct{})

	go connectionCheckHealthy(conn, closeChan)

	for {
		select {
		case <-closeChan:
			fmt.Println("receive close signal")

			if err := conn.Close(); err != nil {
				fmt.Println("connection close error", err.Error())
			}

			break
		default:
			messageType, data, err := conn.ReadMessage()

			fmt.Println("message type is", messageType)

			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseMessage, websocket.CloseAbnormalClosure) {
					fmt.Println("ws connection closed from client")

					if err := conn.Close(); err != nil {
						c.Logger().Error(err)
					}
					fmt.Println("connection closed from server")

					return nil
				}

				if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
					fmt.Println("ws connection unexpect closed from client")

					if err := conn.Close(); err != nil {
						c.Logger().Error(err)
					}
					fmt.Println("connection closed from server")
				}

				fmt.Println("message receive error")

				if err := conn.Close(); err != nil {
					c.Logger().Error(err)
				}
				fmt.Println("connection closed from server")
			}

			fmt.Println("message receive")

			if err := conn.WriteMessage(messageType, data); err != nil {
				c.Logger().Error(err)
			}
		}
	}
}

func connectionCheckHealthy(conn *websocket.Conn, closeChan chan<- struct{}) {
	for {
		if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second*5)); err != nil {
			if websocket.IsCloseError(err, websocket.CloseMessage, websocket.CloseAbnormalClosure) {
				fmt.Println("client closed connection")
			} else {
				fmt.Println("an error occurred")
			}

			closeChan <- struct{}{}

			return
		}

		fmt.Println("connection is ok")
		time.Sleep(time.Second * 1)
	}
}
