package websocket

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

type Config struct {
	Host string
	Port string
}

type Server struct {
	cfg      Config
	log      *slog.Logger
	upgrader websocket.Upgrader
	handler  func(conn *websocket.Conn)
	httpSrv  *http.Server
}

func New(conf Config, log *slog.Logger, handler func(conn *websocket.Conn)) Server {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return Server{
		cfg:      conf,
		log:      log,
		upgrader: upgrader,
		handler:  handler,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)
	s.log.Info("starting WebSocket server at", slog.String("addr", addr))

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			s.log.Error("failed to upgrade to websocket", slog.String("error", err.Error()))
			return
		}
		s.handler(conn)
	})

	s.httpSrv = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.log.Error("websocket server error", slog.String("error", err.Error()))
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("shutting down WebSocket server gracefully...")

	if err := s.httpSrv.Shutdown(ctx); err != nil {
		s.log.Error("error during shutdown", slog.String("error", err.Error()))
		return err
	}

	s.log.Info("WebSocket server shut down cleanly")
	return nil
}
