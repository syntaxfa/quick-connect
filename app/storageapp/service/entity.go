package service

import (
	"time"

	"github.com/syntaxfa/quick-connect/types"
)

type File struct {
	ID         types.ID `json:"id"`
	UploaderID types.ID `json:"uploader_id"`

	Name    string `json:"name"`
	Key     string `json:"key"`
	MimType string `json:"mim_type"`
	Size    int64  `json:"size"`

	Driver   string `json:"driver"`
	Bucket   string `json:"bucket"`
	IsPublic bool   `json:"is_public"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (f File) IsDeleted() bool {
	return f.DeletedAt != nil
}
