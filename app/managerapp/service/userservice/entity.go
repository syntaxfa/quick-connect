package userservice

import (
	"time"

	"github.com/syntaxfa/quick-connect/types"
)

type User struct {
	ID             types.ID     `json:"id"`
	Username       string       `json:"username"`
	HashedPassword string       `json:"-"`
	Fullname       string       `json:"fullname"`
	Email          string       `json:"email"`
	PhoneNumber    string       `json:"phone_number"`
	Avatar         string       `json:"avatar"`
	Roles          []types.Role `json:"roles"`
	LastOnlineAt   time.Time    `json:"last_online_at"`
}
