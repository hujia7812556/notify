package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"notify/internal/config"
)

// TokenManager 处理access token的获取和刷新
type TokenManager struct {
	config     config.WeComConfig
	token      string
	tokenMutex sync.RWMutex
	tokenExp   time.Time
}

func NewTokenManager(config config.WeComConfig) *TokenManager {
	return &TokenManager{
		config: config,
	}
}

func (tm *TokenManager) GetToken(ctx context.Context) (string, error) {
	tm.tokenMutex.RLock()
	if tm.token != "" && time.Now().Before(tm.tokenExp) {
		token := tm.token
		tm.tokenMutex.RUnlock()
		return token, nil
	}
	tm.tokenMutex.RUnlock()

	tm.tokenMutex.Lock()
	defer tm.tokenMutex.Unlock()

	// Double check
	if tm.token != "" && time.Now().Before(tm.tokenExp) {
		return tm.token, nil
	}

	// 企业微信的token获取接口
	url := fmt.Sprintf(
		"https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s",
		tm.config.CorpID,
		tm.config.AppSecret,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("create request failed: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("get access token failed: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response failed: %w", err)
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("get access token failed: %s", result.ErrMsg)
	}

	tm.token = result.AccessToken
	tm.tokenExp = time.Now().Add(time.Duration(result.ExpiresIn-300) * time.Second)

	return tm.token, nil
}
