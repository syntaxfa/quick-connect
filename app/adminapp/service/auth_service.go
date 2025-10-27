package service

import (
	"context"

	"github.com/syntaxfa/quick-connect/app/adminapp/adapter"
)

// AuthService handles authentication logic
type AuthService struct {
	managerClient *adapter.ManagerAppClient
}

// NewAuthService creates a new AuthService
func NewAuthService(managerClient *adapter.ManagerAppClient) *AuthService {
	return &AuthService{
		managerClient: managerClient,
	}
}

// Login calls the manager app client to login
func (s *AuthService) Login(ctx context.Context, username, password string) (*adapter.LoginResponse, error) {
	// در آینده می‌توان منطق بیشتری اینجا اضافه کرد
	// مثلاً بررسی اینکه آیا کاربر نقش ادمین دارد یا خیر
	return s.managerClient.Login(ctx, username, password)
}
