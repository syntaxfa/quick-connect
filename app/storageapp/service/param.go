package service

import (
	"io"

	"github.com/syntaxfa/quick-connect/types"
)

type UploadRequest struct {
	UploaderID  types.ID  `json:"-"`
	File        io.Reader `json:"file"`
	Filename    string    `json:"filename"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	IsPublic    bool      `json:"is_public"`
}
