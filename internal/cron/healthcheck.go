package cron

import (
	"context"
	"fmt"
	"strings"
	"time"

	"notify/internal/config"
	"notify/internal/parser"
	"notify/internal/sender"
	"notify/pkg/logger"

	"go.uber.org/zap"
)

type HealthChecker struct {
	sender *sender.Manager
	config config.HealthCheckConfig
	stop   chan struct{}
}

func NewHealthChecker(sender *sender.Manager, config config.HealthCheckConfig) *HealthChecker {
	return &HealthChecker{
		sender: sender,
		config: config,
		stop:   make(chan struct{}),
	}
}

func (h *HealthChecker) Start() {
	if !h.config.Enabled {
		logger.Info("Health check is disabled")
		return
	}
	go h.run()
}

func (h *HealthChecker) Stop() {
	close(h.stop)
}

func (h *HealthChecker) run() {
	// 解析配置的检查时间
	checkTime := strings.Split(h.config.CheckTime, ":")
	if len(checkTime) != 2 {
		logger.Error("Invalid check time format", zap.String("check_time", h.config.CheckTime))
		return
	}

	hour, min := 8, 0 // 默认值
	fmt.Sscanf(checkTime[0], "%d", &hour)
	fmt.Sscanf(checkTime[1], "%d", &min)

	// 计算下次检查时间
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())
	if now.After(next) {
		next = next.Add(24 * time.Hour)
	}

	logger.Info("Health check scheduled",
		zap.String("next_check", next.Format("2006-01-02 15:04:05")))

	timer := time.NewTimer(time.Until(next))
	defer timer.Stop()

	for {
		select {
		case <-h.stop:
			return
		case <-timer.C:
			// 发送健康检查消息
			h.check()

			// 重置定时器到下一个检查时间
			next = next.Add(24 * time.Hour)
			timer.Reset(time.Until(next))

			logger.Info("Next health check scheduled",
				zap.String("next_check", next.Format("2006-01-02 15:04:05")))
		}
	}
}

func (h *HealthChecker) check() {
	msg := &parser.Message{
		Platform: parser.PlatformWeChat,
		Content:  fmt.Sprintf("系统运行正常\n检查时间: %s", time.Now().Format("2006-01-02 15:04:05")),
		Summary:  "每日健康检查",
		Extra:    make(map[string]any),
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.config.Timeout)
	defer cancel()

	if err := h.sender.Send(ctx, msg); err != nil {
		logger.Error("Health check failed",
			zap.Error(err),
			zap.Time("check_time", time.Now()))
		return
	}

	logger.Info("Health check completed successfully",
		zap.String("check_time", time.Now().Format("2006-01-02 15:04:05")))
}
