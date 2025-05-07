package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"io"
	"time"
)

/*
Control Messages ¶
The WebSocket protocol defines three types of control messages: close, ping and pong. Call the connection WriteControl, WriteMessage or NextWriter methods to send a control message to the peer.

Connections handle received close messages by calling the handler function set with the SetCloseHandler method and by returning a *CloseError from the NextReader, ReadMessage or the message Read method. The default close handler sends a close message to the peer.

Connections handle received ping messages by calling the handler function set with the SetPingHandler method. The default ping handler sends a pong message to the peer.

Connections handle received pong messages by calling the handler function set with the SetPongHandler method. The default pong handler does nothing. If an application sends ping messages, then the application should set a pong handler to receive the corresponding pong.

The control message handler functions are called from the NextReader, ReadMessage and message reader Read methods. The default close and ping handlers can block these methods for a short time when the handler writes to the connection.

The application must read the connection to process close, ping and pong messages sent from the peer. If the application is not otherwise interested in messages from the peer, then the application should start a goroutine to read and discard messages from the peer. A simple example is:

Note:
The Close and WriteControl methods can be called concurrently with all other methods.
*/

func main() {
	httpServer()
}

/*
Buffers ¶
Connections buffer network input and output to reduce the number of system calls when reading or writing messages.

Write buffers are also used for constructing WebSocket frames. See RFC 6455, Section 5 for a discussion of message framing. A WebSocket frame header is written to the network each time a write buffer is flushed to the network. Decreasing the size of the write buffer can increase the amount of framing overhead on the connection.

The buffer sizes in bytes are specified by the ReadBufferSize and WriteBufferSize fields in the Dialer and Upgrader. The Dialer uses a default size of 4096 when a buffer size field is set to zero. The Upgrader reuses buffers created by the HTTP server when a buffer size field is set to zero. The HTTP server buffers have a size of 4096 at the time of this writing.

The buffer sizes do not limit the size of a message that can be read or written by a connection.

Buffers are held for the lifetime of the connection by default. If the Dialer or Upgrader WriteBufferPool field is set, then a connection holds the write buffer only when writing a message.

Applications should tune the buffer sizes to balance memory use and performance. Increasing the buffer size uses more memory, but can reduce the number of system calls to read or write the network. In the case of writing, increasing the buffer size can reduce the number of frame headers written to the network.

Some guidelines for setting buffer parameters are:

Limit the buffer sizes to the maximum expected message size. Buffers larger than the largest message do not provide any benefit.

Depending on the distribution of message sizes, setting the buffer size to a value less than the maximum expected message size can greatly reduce memory use with a small impact on performance. Here's an example: If 99% of the messages are smaller than 256 bytes and the maximum message size is 512 bytes, then a buffer size of 256 bytes will result in 1.01 more system calls than a buffer size of 512 bytes. The memory savings is 50%.

A write buffer pool is useful when the application has a modest number writes over a large number of connections. when buffers are pooled, a larger buffer size has a reduced impact on total memory use and has the benefit of reducing system calls and frame overhead.
*/
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func httpServer() {
	e := echo.New()

	e.GET("/read-loop/", readLoopHandler)

	e.Start(":2525")
}

func readLoopHandler(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err)
	}

	fmt.Println("upgrade http to ws")

	closeChan := make(chan struct{})
	go readLoop(conn, closeChan)

	for {
		select {
		case <-closeChan:
			c.Logger().Error(fmt.Sprintf("connection with ip %s closed.", c.Request().RemoteAddr))
		default:
			messageType, r, err := conn.NextReader()
			if err != nil {
				c.Logger().Error(err)
			}

			fmt.Println("receive a message")

			w, err := conn.NextWriter(messageType)
			if err != nil {
				c.Logger().Error(err)
			}

			if _, err := io.Copy(w, r); err != nil {
				c.Logger().Error(err)
			}

			if err := w.Close(); err != nil {
				c.Logger().Error(err)
			}
		}
	}
}

func readLoop(c *websocket.Conn, closeChan chan<- struct{}) {
	for {
		if err := c.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second*5)); err != nil {
			fmt.Println("connection is closed.")
			c.Close()
			closeChan <- struct{}{}
			break
		}
		time.Sleep(time.Second * 1)
		fmt.Println("connection ok")
	}
}
