package tokenmanager

import (
	"context"
	"sync"
	"time"

	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/authpb"
	"github.com/syntaxfa/quick-connect/types"
)

const defaultTimeAdd = time.Second * 10

type Auth interface {
	Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error)
	TokenRefresh(ctx context.Context, req *authpb.TokenRefreshRequest) (*authpb.TokenRefreshResponse, error)
}

type TokenManager struct {
	mu            sync.RWMutex
	auth          Auth
	accessToken   string
	refreshToken  string
	accessExpiry  time.Time
	refreshExpiry time.Time
	username      string
	password      string
}

func NewTokenManager(username, password string, auth Auth) *TokenManager {
	return &TokenManager{
		auth:     auth,
		username: username,
		password: password,
	}
}

func (t *TokenManager) SetTokenInContext(ctx context.Context) (context.Context, error) {
	t.mu.RLock()
	if t.isAccessTokenValid() {
		defer t.mu.RUnlock()

		return t.contextWithAuthorizationKey(ctx), nil
	}
	t.mu.RUnlock()

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.isAccessTokenValid() {
		return t.contextWithAuthorizationKey(ctx), nil
	}

	if t.isRefreshTokenValid() {
		rErr := t.performRefresh(ctx)
		if rErr == nil {
			return t.contextWithAuthorizationKey(ctx), nil
		}
	}

	if lErr := t.performLogin(ctx); lErr != nil {
		return nil, lErr
	}

	return t.contextWithAuthorizationKey(ctx), nil
}

func (t *TokenManager) contextWithAuthorizationKey(ctx context.Context) context.Context {
	return context.WithValue(ctx, types.AuthorizationKey, t.accessToken)
}

func (t *TokenManager) isAccessTokenValid() bool {
	return t.accessToken != "" && time.Now().Add(defaultTimeAdd).Before(t.accessExpiry)
}

func (t *TokenManager) isRefreshTokenValid() bool {
	return t.refreshToken != "" && time.Now().Add(defaultTimeAdd).Before(t.refreshExpiry)
}

func (t *TokenManager) performLogin(ctx context.Context) error {
	resp, lErr := t.auth.Login(ctx, &authpb.LoginRequest{
		Username: t.username,
		Password: t.password,
	})
	if lErr != nil {
		return lErr
	}

	t.updateTokenFields(resp.GetAccessToken(), resp.GetRefreshToken(), resp.GetAccessExpiresIn(), resp.GetRefreshExpiresIn())

	return nil
}

func (t *TokenManager) performRefresh(ctx context.Context) error {
	resp, rErr := t.auth.TokenRefresh(ctx, &authpb.TokenRefreshRequest{RefreshToken: t.refreshToken})
	if rErr != nil {
		return rErr
	}

	t.updateTokenFields(resp.GetAccessToken(), resp.GetRefreshToken(), resp.GetAccessExpiresIn(), resp.GetRefreshExpiresIn())

	return nil
}

func (t *TokenManager) updateTokenFields(accessToken, refreshToken string, accessExp, refreshExp int32) {
	t.accessToken = accessToken
	t.refreshToken = refreshToken
	t.accessExpiry = time.Now().Add(time.Duration(accessExp) * time.Second)
	t.refreshExpiry = time.Now().Add(time.Duration(refreshExp) * time.Second)
}
