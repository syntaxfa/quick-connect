package userservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/syntaxfa/quick-connect/pkg/cachemanager"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

type UserIDCacheValue struct {
	UserID types.ID `json:"user_id"`
}

func (s Service) GetUserIDFromExternalUserID(ctx context.Context, externalUserID string) (types.ID, error) {
	const op = "service.external_user_id.GetUserIDFromExternalUserID"

	if _, parseErr := ulid.Parse(externalUserID); parseErr == nil {
		return types.ID(externalUserID), nil
	}

	key := s.getUserIDCacheKey(externalUserID)
	var cacheValue UserIDCacheValue
	gErr := s.cache.Get(ctx, key, &cacheValue)

	if gErr == nil {
		return cacheValue.UserID, nil
	}

	if !errors.Is(gErr, cachemanager.ErrKeyNotFound) {
		return "", errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return s.handlerCacheMiss(ctx, externalUserID, key)
}

func (s Service) handlerCacheMiss(ctx context.Context, externalUserID string, cacheKey string) (types.ID, error) {
	const op = "service.external_user_id.handlerCacheMiss"

	exist, existErr := s.externalUserRepo.IsExistUserIDFromExternalUserID(ctx, externalUserID)
	if existErr != nil {
		return "", errlog.ErrContext(ctx, richerror.New(op).WithWrapError(existErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	var userID types.ID

	if exist {
		var guErr error
		userID, guErr = s.externalUserRepo.GetUserIDFromExternalUserID(ctx, externalUserID)
		if guErr != nil {
			return "", errlog.ErrContext(ctx, richerror.New(op).WithWrapError(guErr).WithKind(richerror.KindUnexpected), s.logger)
		}

		return s.setCacheUserID(ctx, cacheKey, userID)
	}

	userID = types.ID(ulid.Make().String())

	if createErr := s.externalUserRepo.CreateUserIDFromExternalUserID(ctx, externalUserID, userID); createErr != nil {
		return "", errlog.ErrContext(ctx, richerror.New(op).WithWrapError(createErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return s.setCacheUserID(ctx, cacheKey, userID)
}

func (s Service) setCacheUserID(ctx context.Context, key string, userID types.ID) (types.ID, error) {
	const op = "service.external_user_id.setCacheUserID"

	if sErr := s.cache.Set(ctx, key, UserIDCacheValue{UserID: userID}, s.cfg.UserIDCacheExpiration); sErr != nil {
		return "", errlog.ErrContext(ctx, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return userID, nil
}

func (s Service) getUserIDCacheKey(externalUserID string) string {
	return fmt.Sprintf("users:externals:%s", externalUserID)
}
