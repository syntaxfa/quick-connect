package httpserver

type Cors struct {
	AllowOrigins []string `koanf:"allow_origins"`
}

type Config struct {
	Port int  `koanf:"port"`
	Cors Cors `koang:"cors"`
}
