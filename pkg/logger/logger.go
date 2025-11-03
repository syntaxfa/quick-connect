package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	FilePath         string `koanf:"file_path"`
	UseLocalTime     bool   `koanf:"use_local_time"`
	FileMaxSizeInMB  int    `koanf:"file_max_size_in_mb"`
	FileMaxAgeInDays int    `koanf:"file_max_age_in_days"`
	MaxBackup        int    `koanf:"max_backup"`
	Compress         bool   `koanf:"compress"`
}

func New(cfg Config, opt *slog.HandlerOptions, writeInConsole bool, serviceName string) *slog.Logger {
	if cfg.FilePath == "" {
		panic("filepath can be blank")
	}

	workingDir, wErr := os.Getwd()
	if wErr != nil {
		panic(fmt.Errorf("error getting current working directory: %w", wErr))
	}

	fileWriter := &lumberjack.Logger{
		Filename:   filepath.Join(workingDir, cfg.FilePath),
		MaxSize:    cfg.FileMaxSizeInMB,
		MaxAge:     cfg.FileMaxAgeInDays,
		MaxBackups: cfg.MaxBackup,
		LocalTime:  cfg.UseLocalTime,
		Compress:   cfg.Compress,
	}

	writers := []io.Writer{fileWriter}
	if writeInConsole {
		writers = append(writers, os.Stdout)
	}

	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(writers...), opt))

	return logger.With("service_name", serviceName)
}
