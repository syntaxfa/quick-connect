package rolemanager

import "github.com/syntaxfa/quick-connect/types"

type RoleManager struct {
	methodRoles map[string][]types.Role
}

func NewRoleManager(methodRoles map[string][]types.Role) *RoleManager {
	return &RoleManager{
		methodRoles: methodRoles,
	}
}

func (r *RoleManager) GetRequireRoles(method string) []types.Role {
	roles, ok := r.methodRoles[method]
	if !ok {
		return []types.Role{types.RoleSuperUser}
	}

	return roles
}
