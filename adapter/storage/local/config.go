package local

type Config struct {
	RootPath string `koanf:"root_path"`
	BaseURL  string `koanf:"base_url"`
}
