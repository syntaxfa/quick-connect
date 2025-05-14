package postgres

import (
	"context"

	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
)

func (d *DB) GetUserByUsername(_ context.Context, _ string) (userservice.User, error) {
	return userservice.User{}, nil
}
