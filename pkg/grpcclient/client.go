package grpcclient

import (
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn *grpc.ClientConn
}

func New(cfg Config) (*Client, error) {
	var opts []grpc.DialOption

	if !cfg.SSLMode {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	if cfg.UseOtel {
		opts = append(opts, grpc.WithStatsHandler(otelgrpc.NewClientHandler(
			otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
			otelgrpc.WithMeterProvider(otel.GetMeterProvider()),
		)))
	}

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), opts...)
	if err != nil {
		return &Client{}, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Conn() *grpc.ClientConn {
	return c.conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}
