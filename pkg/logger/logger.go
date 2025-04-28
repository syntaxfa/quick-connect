package logger

import (
	"io"
	"log"
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

var globalLogger *slog.Logger

func L() *slog.Logger {
	return globalLogger
}

func SetDefault(cfg Config, opt *slog.HandlerOptions, writeInConsole bool) {
	workingDir, wErr := os.Getwd()
	if wErr != nil {
		log.Fatalf("error getting current working directory, %s", wErr.Error())
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

	globalLogger = slog.New(slog.NewJSONHandler(io.MultiWriter(writers...), opt))
}

func New(cfg Config, opt *slog.HandlerOptions, writeInConsole bool) *slog.Logger {
	workingDir, wErr := os.Getwd()
	if wErr != nil {
		log.Fatalf("error getting current working directory, %s", wErr.Error())
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

	return slog.New(slog.NewJSONHandler(io.MultiWriter(writers...), opt))
}
