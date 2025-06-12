package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type GorillaConnection struct {
	conn *websocket.Conn
}

func (c *GorillaConnection) ReadMessage() (messageType int, p []byte, err error) {
	return c.conn.ReadMessage()
}

func (c *GorillaConnection) WriteMessage(messageType int, data []byte) error {
	return c.conn.WriteMessage(messageType, data)
}

func (c *GorillaConnection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *GorillaConnection) Close() error {
	return c.conn.Close()
}

type GorillaUpgrader struct {
	upgrader websocket.Upgrader
}

func NewGorillaUpgrader(cfg Config, checkOrigin func(r *http.Request) bool) *GorillaUpgrader {
	return &GorillaUpgrader{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  cfg.ReadBufferSize,
			WriteBufferSize: cfg.WriteBufferSize,
			CheckOrigin:     checkOrigin,
		},
	}
}

func (u *GorillaUpgrader) Upgrade(w http.ResponseWriter, r *http.Request) (*GorillaConnection, error) {
	conn, err := u.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return &GorillaConnection{conn: conn}, nil
}
