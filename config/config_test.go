package config_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/syntaxfa/quick-connect/config"
)

type application struct {
	Name string `koanf:"name"`
	Port int    `koanf:"port"`
}

type dbOptions struct {
	MaxConnSize    int `koanf:"max_conn_size"`
	MaxThreadCount int `koanf:"max_thread_count"`
}

type db struct {
	User     string    `koanf:"user"`
	Password string    `koanf:"password"`
	Host     string    `koanf:"host"`
	Port     int       `koanf:"port"`
	Options  dbOptions `koanf:"options"`
}

type Config struct {
	Debug       bool        `koanf:"debug"`
	Application application `koanf:"application"`
	DB          db          `koanf:"db"`
}

var options = config.Option{
	Prefix:       "QUICK_",
	Delimiter:    ".",
	Separator:    "__",
	YamlFilePath: "",
}

func getYamlFilePath() string {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("an error occurred when trying get current directory, err: %s\n", err.Error())
	}

	return filepath.Join(workingDir, "config.yml")
}

func defaultConfig() Config {
	return Config{
		Debug: true,
		Application: application{
			Name: "my app",
			Port: 8080,
		},
		DB: db{
			Host:    "localhost",
			Options: dbOptions{MaxThreadCount: 22},
		},
	}
}

func TestLoadingDefaultConfigFromStruct(t *testing.T) {
	var cfg Config

	expectedStruct := Config{
		Debug: true,
		Application: application{
			Name: "my app",
			Port: 8080,
		},
		DB: db{
			Host:    "localhost",
			Options: dbOptions{MaxThreadCount: 22},
		},
	}

	config.Load(options, &cfg, defaultConfig())

	if !reflect.DeepEqual(expectedStruct, cfg) {
		t.Fatalf("expected: %+v \ngot: %+v\n", expectedStruct, cfg)
	}
}

func TestLoadingConfigFromYMLFile(t *testing.T) {
	var cfg Config

	options.YamlFilePath = getYamlFilePath()

	expectedCfg := Config{
		Application: application{
			Name: "my app",
			Port: 8080,
		},
		Debug: false,
		DB: db{
			Host: "postgres.quick.club",
			Options: dbOptions{
				MaxThreadCount: 35,
			},
		},
	}

	ymlConfig := []byte(`debug: false
db:
  host: postgres.quick.club
  options:
    max_thread_count: 35
`)

	ymlFile, _ := os.Create("config.yml")
	defer func() {
		if err := ymlFile.Close(); err != nil {
			fmt.Printf("%s", err.Error())
		}

		if err := os.Remove("config.yml"); err != nil {
			fmt.Printf("%s", err.Error())
		}
	}()

	if _, err := ymlFile.Write(ymlConfig); err != nil {
		fmt.Printf("%s", err.Error())
	}

	config.Load(options, &cfg, defaultConfig())

	if !reflect.DeepEqual(cfg, expectedCfg) {
		t.Fatalf("expected: %+v \ngot: %+v\n", expectedCfg, cfg)
	}
}

func TestLoadingConfigFromEnvironment(t *testing.T) {
	options.YamlFilePath = ""

	var cfg Config

	expectedStruct := Config{
		Debug: true,
		Application: application{
			Name: "my app",
			Port: 8080,
		},
		DB: db{
			Host:    "localhost",
			Options: dbOptions{MaxThreadCount: 65},
		},
	}

	t.Setenv(fmt.Sprintf("%s%s__%s__%s", options.Prefix, "db", "options", "max_thread_count"), "65")
	defer os.Unsetenv(fmt.Sprintf("%s%s__%s__%s", options.Prefix, "db", "options", "max_thread_count"))

	config.Load(options, &cfg, defaultConfig())

	if !reflect.DeepEqual(expectedStruct, cfg) {
		t.Fatalf("expected: %+v \ngot: %+v\n", expectedStruct, cfg)
	}
}
