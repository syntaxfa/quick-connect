package websocket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type GorillaConnection struct {
	conn *websocket.Conn
}

func (c *GorillaConnection) ReadMessage() (int, []byte, error) {
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

// SetReadLimit sets the maximum size for a message read from the peer.
func (c *GorillaConnection) SetReadLimit(limit int64) {
	c.conn.SetReadLimit(limit)
}

// SetReadDeadline sets the read deadline on the underlying connection.
func (c *GorillaConnection) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the write deadline on the underlying connection.
func (c *GorillaConnection) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

// SetPongHandler sets the handler for pong messages received from the peer.
func (c *GorillaConnection) SetPongHandler(h func(appData string) error) {
	c.conn.SetPongHandler(h)
}

// SetCloseHandler sets the handler for close messages received from the peer.
func (c *GorillaConnection) SetCloseHandler(h func(code int, text string) error) {
	c.conn.SetCloseHandler(h)
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

func (u *GorillaUpgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*GorillaConnection, error) {
	conn, err := u.upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}

	return &GorillaConnection{conn: conn}, nil
}

// We define a dummy type here just for the static check.
type serviceConnection interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
	SetReadLimit(limit int64)
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	SetPongHandler(h func(appData string) error)
	SetCloseHandler(h func(code int, text string) error)
}

var _ serviceConnection = (*GorillaConnection)(nil)
