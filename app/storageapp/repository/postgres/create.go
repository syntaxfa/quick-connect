package postgres

import (
	"context"

	"github.com/syntaxfa/quick-connect/app/storageapp/service"
)

func (d *DB) Save(_ context.Context, _ service.File) error {
	return nil
}
