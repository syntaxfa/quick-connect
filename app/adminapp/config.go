package adminapp

import "github.com/syntaxfa/quick-connect/pkg/httpserver"

type Config struct {
	HTTPServer httpserver.Config `koanf:"http_server"`
}
