package service

import (
	"time"

	"github.com/syntaxfa/quick-connect/types"
)

type User struct {
	ID             types.ID  `json:"id"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"-"`
	Fullname       string    `json:"fullname"`
	Avatar         string    `json:"avatar"`
	LastOnlineAt   time.Time `json:"last_online_at"`
	Role           UserRole  `json:"role"`
}

type UserRole uint8

const (
	UserRoleSuperUser = iota + 1
	UserRoleAdmin
)
