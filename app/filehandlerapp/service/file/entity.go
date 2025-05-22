package file

import (
	"fmt"
	"time"

	"github.com/syntaxfa/quick-connect/types"
)

type (
	FileType    string
	StorageType string
)

const (
	// File types
	FileTypeChat FileType = "chat"

	// File storage
	FileStorageTypeLocal StorageType = "local"
	FileStorageTypeS3    StorageType = "s3"
)

type File struct {
	ID          types.ULID
	Type        FileType
	TypeID      types.ID
	Extension   string
	StorageType StorageType
	Size        int // in megabytes
	ContentType string
	IsPublic    bool
	IsDeleted   bool
	CreatedAt   time.Time
}

func (f *File) GetFileName() string {
	return fmt.Sprintf("%v.%v", f.ID, f.Extension)
}
