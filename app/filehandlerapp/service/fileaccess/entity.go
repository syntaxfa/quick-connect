package fileaccess

import (
	"github.com/syntaxfa/quick-connect/app/filehandlerapp/service/file"
	"github.com/syntaxfa/quick-connect/types"
)

type (
	PermissionType string
)

const (
	// Permissions
	PermissionsRead  PermissionType = "r"
	PermissionWrite  PermissionType = "w"
	PermissionUpdate PermissionType = "u"
	PermissionDelete PermissionType = "d"
)

type FileAccess struct {
	ClientID    types.ID
	File        file.File
	Permissions []PermissionType
}
