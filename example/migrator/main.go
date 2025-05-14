package main

import (
	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/pkg/migrator"
	"log"
)

func main() {
	cfg := postgres.Config{
		Host:            "localhost",
		Port:            5432,
		Username:        "test",
		Password:        "test",
		DBName:          "test",
		SSLMode:         "disable",
		MaxIdleConns:    5,
		MaxOpenConns:    10,
		ConnMaxLifetime: 600,
	}

	migrate := migrator.New(migrator.Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		DBName:   cfg.DBName,
		SSLMode:  cfg.SSLMode,
	}, "example/migrator/migrations")

	if n, uErr := migrate.Down(1); uErr != nil {
		log.Fatalln(uErr.Error())
	} else {
		log.Printf("migrations %d\n", n)
	}
}
