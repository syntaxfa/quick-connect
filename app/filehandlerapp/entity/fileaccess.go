package entity

import "github.com/syntaxfa/quick-connect/types"

const (
	PermissionsRead  = "r"
	PermissionWrite  = "w"
	PermissionUpdate = "u"
	PermissionDelete = "d"
)

type FileAccess struct {
	ClientID    types.ID
	File        types.ID
	Permissions []string
}
