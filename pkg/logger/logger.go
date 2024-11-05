package logger

import (
	"notify/pkg/types"

	"go.uber.org/zap"
)

var log *zap.Logger

func Init(mode string, cfg types.LogConfig) error {
	config := zap.NewProductionConfig()

	// 根据运行模式选择输出
	if mode == "debug" {
		if cfg.Debug.Output != "" {
			config.OutputPaths = []string{cfg.Debug.Output}
		} else {
			config.OutputPaths = []string{"stdout"}
		}
	} else {
		if cfg.Release.Output != "" {
			config.OutputPaths = []string{cfg.Release.Output}
		} else {
			config.OutputPaths = []string{cfg.Output}
		}
	}

	// 设置日志级别
	switch cfg.Level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	// 设置格式
	if cfg.Format == "json" {
		config.Encoding = "json"
	} else {
		config.Encoding = "console"
	}

	var err error
	log, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}
