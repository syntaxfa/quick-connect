package postgres

import (
	"context"

	"github.com/syntaxfa/quick-connect/types"
)

func (d *DB) MarkAllAsReadByUserID(_ context.Context, _ types.ID) error {
	// TODO: implemented
	return nil
}

func (d *DB) MarkAsRead(_ context.Context, _ types.ID) error {
	// TODO: implemented
	return nil
}
