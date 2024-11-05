package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"notify/internal/config"
	"notify/internal/dispatcher"
	"notify/internal/parser"
	"notify/internal/sender"
	"notify/internal/server"
	"notify/pkg/logger"

	"notify/internal/sender/factory"
	"notify/internal/sender/wechat"

	"notify/internal/cron"

	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	if err := logger.Init(cfg.Server.Mode, cfg.Log); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// 初始化发送器
	senderMgr := sender.NewManager()

	// 创建并注册微信发送器
	wechatSender, err := factory.CreateWeChatSender(wechat.WeChatSenderType(cfg.WeChat.SenderType), cfg.WeChat)
	if err != nil {
		log.Fatalf("Failed to create WeChat sender: %v", err)
	}
	senderMgr.Register(parser.PlatformWeChat, wechatSender)

	// 注册钉钉发送器
	senderMgr.Register(parser.PlatformDingTalk, sender.NewDingTalkSender(cfg.DingTalk))

	// 初始化分发器
	disp := dispatcher.New(cfg.Dispatcher.BufferSize, cfg.Dispatcher.WorkerPoolSize, senderMgr)

	// 启动分发器
	ctx, cancel := context.WithCancel(context.Background())
	disp.Start(ctx)

	// 初始化并启动服务器
	srv := server.New(cfg.Server, disp)
	go func() {
		if err := srv.Start(); err != nil {
			logger.Error("Server error", zap.Error(err))
			cancel()
		}
	}()

	// 初始化健康检查器
	healthChecker := cron.NewHealthChecker(senderMgr, cfg.HealthCheck)
	healthChecker.Start()

	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// 优雅关闭
	healthChecker.Stop()
	cancel()
	disp.Stop()
	logger.Info("Server shutdown complete")
}
