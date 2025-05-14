package userservice

import (
	"time"

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
