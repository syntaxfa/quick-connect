package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	writeWait      = 10 * time.Second    // waiting tome for writing
	pongWait       = 60 * time.Second    // waiting time for receive pong from client
	pingPeriod     = (pongWait * 9) / 10 // waiting time for sending ping to client(less than pongWait)
	maxMessageSize = 1024                // maximum size of message
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(_ *http.Request) bool {
			// check origins
			return true
		},
	}

	// connection management
	connections = struct {
		sync.RWMutex
		clients map[*websocket.Conn]bool
	}{
		clients: make(map[*websocket.Conn]bool),
	}
)

type Client struct {
	conn      *websocket.Conn
	send      chan []byte
	ctx       context.Context
	cancelCtx context.CancelFunc
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/messages", messagesHandler)

	log.Fatal(e.Start(":2525"))
}

func messagesHandler(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)

		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		conn:      conn,
		send:      make(chan []byte, 256),
		ctx:       ctx,
		cancelCtx: cancel,
	}

	connections.Lock()
	connections.clients[conn] = true
	connections.Unlock()

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))

		return nil
	})

	go client.writePump()
	go client.readPump()

	return nil
}

func (c *Client) readPump() {
	defer func() {
		c.cancelCtx()
		c.conn.Close()

		connections.Lock()
		delete(connections.clients, c.conn)
		connections.Unlock()

		close(c.send)
	}()

	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error: %v", err)
			} else {
				log.Printf("Read error: %v", err)
			}
			break
		}

		log.Printf("Received message type: %d, data: %s", messageType, message)

		c.send <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})

				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("Error getting writer: %v", err)
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				log.Printf("Error closing writer: %v", err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error sending ping: %v", err)

				return
			}
			log.Printf("ping pong")
		case <-c.ctx.Done():
			return
		}
	}
}

//func broadcast(message []byte) {
//	connections.RLock()
//	defer connections.RUnlock()
//
//	for conn := range connections.clients {
//		select {
//		case client.send <- message:
//		default:
//			conn.Close()
//			connections.Lock()
//			delete(connections.clients, conn)
//			connections.Unlock()
//		}
//	}
//}
