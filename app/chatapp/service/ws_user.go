package service

import (
	"encoding/json"
	"log/slog"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

// wsUser holds all common fields and methods for a websocket participant.
type wsUser struct {
	id           string
	hub          *Hub
	conn         Connection
	send         chan Message
	username     string
	logger       *slog.Logger
	cfg          Config
	wsConnection *wsConnection
}

// newWsUser is the common constructor for any websocket participant.
func newWsUser(cfg Config, id string, hub *Hub, conn Connection, username string, logger *slog.Logger, channelSize int) wsUser {
	user := wsUser{
		cfg:      cfg,
		id:       id,
		hub:      hub,
		conn:     conn,
		send:     make(chan Message, channelSize), // Use the provided channel size
		username: username,
		logger:   logger,
	}
	user.wsConnection = newWsConnection(cfg, conn, user.send, logger)
	return user
}

// GetID returns the user's ID.
func (u *wsUser) GetID() string {
	return u.id
}

// Send queues a message to be sent to the user.
func (u *wsUser) Send(message Message) {
	u.send <- message
}

// Close closes the send channel.
func (u *wsUser) Close() {
	close(u.send)
}

// writePump delegates the write operation to the wsConnection.
func (u *wsUser) writePump(op string) {
	u.wsConnection.writePump(op)
}

// readPump is the generic message-reading loop.
// It takes the operation string, the destination broadcast channel,
// and a function to call when unregistering.
func (u *wsUser) readPump(op string, broadcastChan chan<- Message, unregisterFunc func()) {
	defer func() {
		unregisterFunc() // Call the specific unregister logic
		if cErr := u.conn.Close(); cErr != nil {
			errlog.WithoutErr(richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), u.logger)
		}
	}()

	for {
		_, message, rErr := u.conn.ReadMessage()
		if rErr != nil {
			errlog.WithoutErr(richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected), u.logger)
			break
		}

		var msg Message
		if uErr := json.Unmarshal(message, &msg); uErr != nil {
			errlog.WithoutErr(richerror.New(op).WithWrapError(uErr).WithKind(richerror.KindUnexpected), u.logger)
			continue
		}

		msg.Sender = u.id

		// Send to the destination channel provided as an argument
		broadcastChan <- msg
	}
}
