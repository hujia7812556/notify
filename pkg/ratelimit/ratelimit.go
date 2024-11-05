package ratelimit

import (
	"context"
	"time"
)

// RateLimiter 限流器接口
type RateLimiter interface {
	Wait(ctx context.Context) error
}

// TokenBucket 令牌桶限流器
type TokenBucket struct {
	tokens   chan struct{}
	interval time.Duration
	stopChan chan struct{}
}

// NewTokenBucket 创建新的令牌桶限流器
func NewTokenBucket(qps int) *TokenBucket {
	if qps <= 0 {
		qps = 1
	}
	tb := &TokenBucket{
		tokens:   make(chan struct{}, qps),
		interval: time.Second / time.Duration(qps),
		stopChan: make(chan struct{}),
	}

	// 初始填充令牌
	for i := 0; i < qps; i++ {
		tb.tokens <- struct{}{}
	}

	// 启动令牌生成器
	go tb.generateTokens()

	return tb
}

// Wait 等待获取令牌
func (tb *TokenBucket) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-tb.tokens:
		return nil
	}
}

// Stop 停止令牌生成
func (tb *TokenBucket) Stop() {
	close(tb.stopChan)
}

// generateTokens 持续生成令牌
func (tb *TokenBucket) generateTokens() {
	ticker := time.NewTicker(tb.interval)
	defer ticker.Stop()

	for {
		select {
		case <-tb.stopChan:
			return
		case <-ticker.C:
			select {
			case tb.tokens <- struct{}{}:
			default:
				// 令牌桶已满，丢弃令牌
			}
		}
	}
}
