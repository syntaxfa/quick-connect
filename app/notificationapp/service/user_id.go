package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/syntaxfa/quick-connect/pkg/cachemanager"
	"github.com/syntaxfa/quick-connect/types"
)

type UserIDCacheValue struct {
	UserID types.ID `json:"user_id"`
}

func (s Service) getUserIDFromExternalUserID(ctx context.Context, externalUserID string) (types.ID, error) {
	if _, pErr := ulid.Parse(externalUserID); pErr == nil {
		return types.ID(externalUserID), nil
	}

	key := s.getUserIDCacheKey(externalUserID)
	var cacheValue UserIDCacheValue
	gErr := s.cache.Get(ctx, key, &cacheValue)

	if gErr == nil {
		return cacheValue.UserID, nil
	}

	if !errors.Is(gErr, cachemanager.ErrKeyNotFound) {
		return "", gErr
	}

	return s.handleCacheMiss(ctx, externalUserID, key)
}

func (s Service) handleCacheMiss(ctx context.Context, externalUserID string, cacheKey string) (types.ID, error) {
	exist, eErr := s.repo.IsExistUserIDFromExternalUserID(ctx, externalUserID)
	if eErr != nil {
		return "", eErr
	}

	var userID types.ID

	if exist {
		var guErr error
		userID, guErr = s.repo.GetUserIDFromExternalUserID(ctx, externalUserID)
		if guErr != nil {
			return "", guErr
		}

		if setErr := s.setCacheUserID(ctx, cacheKey, userID); setErr != nil {
			return "", setErr
		}

		return userID, nil
	}

	userID = types.ID(ulid.Make().String())

	if setErr := s.repo.CreateUserIDFromExternalUserID(ctx, externalUserID, userID); setErr != nil {
		return "", setErr
	}

	if err := s.setCacheUserID(ctx, cacheKey, userID); err != nil {
		return "", err
	}

	return userID, nil
}

func (s Service) setCacheUserID(ctx context.Context, key string, userID types.ID) error {
	return s.cache.Set(ctx, key, UserIDCacheValue{UserID: userID}, s.cfg.UserIDCacheExpiration)
}

func (s Service) getUserIDCacheKey(externalUserID string) string {
	return fmt.Sprintf("users:externals:%s", externalUserID)
}
