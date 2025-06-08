package entity

import (
	"fmt"
	"time"

	"github.com/syntaxfa/quick-connect/types"
)

const (
	FileTypeChat = "chat"

	FileStorageTypeLocal = "local"
	FileStorageTypeS3    = "s3"
)

type File struct {
	ID          types.ID
	Type        string
	TypeID      types.ID
	Uploaded    time.Time
	Extension   string
	StorageType string
	Size        int // in megabytes
	ContentType string
	IsPublic    bool
	IsDeleted   bool
}

func (f *File) GetFileName() string {
	return fmt.Sprintf("%v.%v", f.ID, f.Extension)
}
