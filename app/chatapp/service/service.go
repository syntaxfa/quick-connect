package service

import (
	"github.com/gorilla/websocket"
	"log/slog"
)

type Service struct {
	logger   *slog.Logger
	clients  map[string]*websocket.Conn
	supports map[string]*websocket.Conn
}

func New(logger *slog.Logger) Service {
	return Service{
		logger:   logger,
		clients:  make(map[string]*websocket.Conn),
		supports: make(map[string]*websocket.Conn),
	}
}
