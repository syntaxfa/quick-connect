package local

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/syntaxfa/quick-connect/app/storageapp/service"
)

var _ service.Storage = (*Adapter)(nil)

const (
	defaultDirPerm = 0750 // Owner: RWX, Group: R-X, Other: ---
)

type Adapter struct {
	cfg    Config
	logger *slog.Logger
}

func New(cfg Config, logger *slog.Logger) (*Adapter, error) {
	if mErr := os.MkdirAll(cfg.RootPath, defaultDirPerm); mErr != nil {
		return nil, fmt.Errorf("failed to create local storage root path %s: %s", cfg.RootPath, mErr.Error())
	}

	return &Adapter{
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (a *Adapter) Upload(_ context.Context, file io.Reader, _ int64, key, _ string, _ bool) (string, error) {
	fullPath := filepath.Join(a.cfg.RootPath, key)

	dir := filepath.Dir(fullPath)
	if mErr := os.MkdirAll(dir, defaultDirPerm); mErr != nil {
		return "", fmt.Errorf("failed to create directory structure for file with dir %s, error: %s", dir, mErr.Error())
	}

	dst, createErr := os.Create(fullPath)
	if createErr != nil {
		return "", fmt.Errorf("failed to create file on disk: %s", createErr.Error())
	}
	defer func() {
		if cErr := dst.Close(); cErr != nil {
			a.logger.Error("destination file can't close", slog.String("error", cErr.Error()))
		}
	}()

	if _, copyErr := io.Copy(dst, file); copyErr != nil {
		return "", fmt.Errorf("failed to write file content: %s", copyErr.Error())
	}

	return key, nil
}

func (a *Adapter) Delete(_ context.Context, key string) error {
	fullPath := filepath.Join(a.cfg.RootPath, key)

	if rErr := os.Remove(fullPath); rErr != nil {
		if os.IsNotExist(rErr) {
			return nil
		}
		return fmt.Errorf("failed to delete local file: %s", rErr.Error())
	}

	return nil
}

func (a *Adapter) GetURL(_ context.Context, key string) (string, error) {
	return fmt.Sprintf("%s/%s", a.cfg.BaseURL, key), nil
}

func (a *Adapter) GetPresignedURL(ctx context.Context, key string) (string, error) {
	return a.GetURL(ctx, key)
}

func (a *Adapter) Exists(_ context.Context, key string) (bool, error) {
	fullPath := filepath.Join(a.cfg.RootPath, key)

	_, sErr := os.Stat(fullPath)
	if sErr == nil {
		return true, nil
	}

	if os.IsNotExist(sErr) {
		return false, nil
	}

	return false, sErr
}
