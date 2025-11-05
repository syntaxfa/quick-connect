package main

import (
	"github.com/syntaxfa/quick-connect/example/grpc/interval/managerauth"
	"github.com/syntaxfa/quick-connect/example/grpc/interval/manageruser"
)

func main() {
	managerauth.ManagerAuth()

	manageruser.ManagerUser()
}
