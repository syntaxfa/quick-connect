package userservice

import (
	"fmt"
	"time"

	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

type User struct {
	ID             types.ID   `json:"id"`
	Username       string     `json:"username"`
	HashedPassword string     `json:"-"`
	Fullname       string     `json:"fullname"`
	Avatar         string     `json:"avatar"`
	Role           types.Role `json:"role"`
	LastOnlineAt   time.Time  `json:"last_online_at"`
}

func RoleStringToInt(roleString string) (types.Role, error) {
	const op = "service.user.entity.RoleStringToInt"

	switch roleString {
	case "superuser":
		return types.RoleSuperUser, nil
	case "admin":
		return types.RoleAdmin, nil
	}

	return 0, richerror.New(op).WithMessage(fmt.Sprintf("unknown role, %s", roleString)).WithKind(richerror.KindUnexpected)
}

func RoleIntToString(role types.Role) (string, error) {
	const op = "service.user.entity.RoleIntToString"

	switch role {
	case types.RoleSuperUser:
		return "superuser", nil
	case types.RoleAdmin:
		return "admin", nil
	}

	return "", richerror.New(op).WithMessage(fmt.Sprintf("unknown role, %d", role))
}
