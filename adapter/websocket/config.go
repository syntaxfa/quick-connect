package websocket

type Config struct {
	ReadBufferSize  int `koanf:"read_buffer_size"`
	WriteBufferSize int `koanf:"write_buffer_size"`
}
