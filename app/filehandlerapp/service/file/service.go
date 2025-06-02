package file

import (
	"context"

	"github.com/syntaxfa/quick-connect/types"
)

type Repository interface {
	Create(ctx context.Context, file File) (types.ULID, error)
	Get(ctx context.Context, ID types.ULID) (File, error)
	Delete(ctx context.Context, ID types.ULID) error
}
