package wechat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"notify/internal/config"
	"notify/pkg/logger"

	"go.uber.org/zap"
)

type WeComSender struct {
	tokenManager *TokenManager
	config       config.WeComConfig
	client       *http.Client
}

func NewWeComSender(config config.WeComConfig) *WeComSender {
	return &WeComSender{
		tokenManager: NewTokenManager(config),
		config:       config,
		client:       &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *WeComSender) Type() WeChatSenderType {
	return SenderTypeWeCom
}

func (s *WeComSender) Send(ctx context.Context, content string, summary string, extra map[string]any) error {
	token, err := s.tokenManager.GetToken(ctx)
	if err != nil {
		return fmt.Errorf("get token failed: %w", err)
	}

	messageContent := content
	if summary != "" {
		messageContent = fmt.Sprintf("【%s】\n\n%s", summary, content)
	}

	msg := struct {
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		AgentID string `json:"agentid"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
	}{
		ToUser:  extra["user_id"].(string),
		MsgType: "text",
		AgentID: s.config.AgentID,
		Text: struct {
			Content string `json:"content"`
		}{
			Content: messageContent,
		},
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message failed: %w", err)
	}

	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response failed: %w", err)
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("send message failed: %s", result.ErrMsg)
	}

	logger.Info("WeCom message sent successfully",
		zap.String("user_id", msg.ToUser))

	return nil
}
