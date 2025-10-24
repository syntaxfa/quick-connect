package main

import (
	"context"
	"fmt"

	"github.com/syntaxfa/quick-connect/adapter/manager"
	"github.com/syntaxfa/quick-connect/pkg/grpcclient"
)

func main() {
	cfg := grpcclient.Config{
		Host:    "localhost",
		Port:    2541,
		SSLMode: false,
		UseOtel: true,
	}

	grpcClient, gErr := grpcclient.New(cfg)
	if gErr != nil {
		panic(gErr)
	}

	tokenAdapter := manager.NewTokenAdapter(grpcClient.Conn())
	fmt.Println(tokenAdapter.GetPublicKey(context.Background()))
}
