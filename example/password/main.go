package main

import (
	"fmt"

	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/randomly"
)

func main() {
	password, gErr := randomly.GeneratePassword(64)
	if gErr != nil {
		panic(gErr)
	}

	fmt.Println(password)

	hash, hErr := userservice.HashPassword(password)
	if hErr != nil {
		panic(hErr)
	}

	fmt.Println(hash)
}
