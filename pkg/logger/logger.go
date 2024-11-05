package logger

import (
	"go.uber.org/zap"
)

// LogConfig 日志配置结构体
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

var log *zap.Logger

func Init(cfg LogConfig) error {
	config := zap.NewProductionConfig()

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

	// 设置输出
	if cfg.Output != "" {
		config.OutputPaths = []string{cfg.Output}
	} else {
		config.OutputPaths = []string{"stdout"}
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
