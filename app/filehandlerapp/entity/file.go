package entity

import (
	"time"

	"github.com/syntaxfa/quick-connect/types"
)

const (
	FileTypeChat = "chat"

	StorageTypeLocal = "local"
	StorageTypeS3    = "s3"
)

type File struct {
	Id          types.ID
	Type        string
	TypeId      types.ID
	Uploaded    time.Time
	Name        string
	StorageType string
}
