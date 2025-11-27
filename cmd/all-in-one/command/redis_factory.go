package command

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/syntaxfa/quick-connect/adapter/redis"
)

type redisFactory struct {
	conns        map[string]*redis.Adapter
	mu           sync.Mutex
	globalLogger *slog.Logger
}

func newRedisFactory(globalLogger *slog.Logger) *redisFactory {
	return &redisFactory{
		conns:        make(map[string]*redis.Adapter),
		globalLogger: globalLogger,
	}
}

// newConnection if connection exists, returns.
func (r *redisFactory) newConnection(cfg redis.Config) *redis.Adapter {
	r.mu.Lock()
	defer r.mu.Unlock()

	conn, ok := r.conns[fmt.Sprintf("%d@%s:%d/%s", cfg.DB, cfg.Host, cfg.Port, cfg.Password)]
	if !ok {
		conn = redis.New(cfg, r.globalLogger)

		r.globalLogger.Info("create redis connection")
	} else {
		r.globalLogger.Info("using same redis connection")
	}

	return conn
}

func (r *redisFactory) closeAll() {
	for key, conn := range r.conns {
		if cErr := conn.Close(); cErr != nil {
			r.globalLogger.Error("redis connection closed failed", slog.String("error", cErr.Error()),
				slog.String("connection", key))
		}
	}

	r.globalLogger.Info("redis connections gracefully shutdown")
}
