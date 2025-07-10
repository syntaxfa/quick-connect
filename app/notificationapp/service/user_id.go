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
	_, pErr := ulid.Parse(externalUserID)
	if pErr == nil {
		return types.ID(externalUserID), nil
	}

	key := s.getUserIDCacheKey(externalUserID)
	var cacheValue UserIDCacheValue
	if gErr := s.cache.Get(ctx, key, &cacheValue); gErr != nil {
		if !errors.Is(gErr, cachemanager.ErrKeyNotFound) {
			return "", gErr
		}

		exist, eErr := s.repo.IsExistUserIDFromExternalUserID(ctx, externalUserID)
		if eErr != nil {
			return "", eErr
		}

		if exist {
			userID, gErr := s.repo.GetUserIDFromExternalUserID(ctx, externalUserID)
			if gErr != nil {
				return "", gErr
			}

			fmt.Println(s.cfg.UserIDCacheExpiration)

			if sErr := s.cache.Set(ctx, key, UserIDCacheValue{UserID: userID}, s.cfg.UserIDCacheExpiration); sErr != nil {
				return "", sErr
			}

			return userID, nil
		}
		userID := types.ID(ulid.Make().String())
		if cErr := s.repo.CreateUserIDFromExternalUserID(ctx, externalUserID, userID); cErr != nil {
			return "", cErr
		}

		if sErr := s.cache.Set(ctx, key, UserIDCacheValue{UserID: userID}, s.cfg.UserIDCacheExpiration); sErr != nil {
			return "", sErr
		}

		return userID, nil
	}

	return cacheValue.UserID, nil
}

func (s Service) getUserIDCacheKey(externalUserID string) string {
	return fmt.Sprintf("users:externals:%s", externalUserID)
}
