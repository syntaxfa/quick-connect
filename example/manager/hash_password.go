package main

import (
	"fmt"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
)

func main() {
	fmt.Println(userservice.HashPassword("Password"))
}
