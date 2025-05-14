package postgres

import "context"

func (d *DB) IsExistUserByUsername(_ context.Context, _ string) (bool, error) {
	return false, nil
}
