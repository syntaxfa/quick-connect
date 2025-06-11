package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"log/slog"
	"os"

	"github.com/oklog/ulid/v2"
	"github.com/syntaxfa/quick-connect/adapter/postgres"
	postgres2 "github.com/syntaxfa/quick-connect/app/notificationapp/repository/postgres"
	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
	"github.com/syntaxfa/quick-connect/types"
)

func main() {
	cfg := postgres.Config{
		Host:            "localhost",
		Port:            11580,
		Username:        "LoPgYJqYGZ53",
		Password:        "8SHDSgdihmMH9EQsXfRZzLHes3F3kgxa",
		DBName:          "defaultdb",
		SSLMode:         "disable",
		MaxIdleConns:    10,
		MaxOpenConns:    20,
		ConnMaxLifetime: 600,
		PathOfMigration: "app/notificationapp/repository/migrations",
	}

	psAdapter := postgres.New(cfg, slog.Default())

	ctx := context.Background()
	repo := postgres2.New(psAdapter)

	data := map[string]string{
		"service": "account",
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	channelDeliveries := []service.ChannelDeliveryRequest{
		{Channel: service.ChannelTypeEmail},
		{Channel: service.ChannelTypeWebPush},
	}

	notification, err := repo.Save(ctx, service.SendNotificationRequest{
		ID:                types.ID(ulid.Make().String()),
		UserID:            types.ID(ulid.Make().String()),
		Type:              service.NotificationTypeInfo,
		Title:             "message title 2",
		Body:              "message body 2",
		Data:              dataJson,
		ChannelDeliveries: channelDeliveries,
	})
	if err != nil {
		errlog.WithoutErr(err, slog.Default())

		os.Exit(1)
	}

	fmt.Printf("%+v", notification)
}
