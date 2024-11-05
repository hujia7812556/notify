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

type WxPusherSender struct {
	config config.WxPusherConfig
	client *http.Client
}

func NewWxPusherSender(config config.WxPusherConfig) *WxPusherSender {
	return &WxPusherSender{
		config: config,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *WxPusherSender) Type() WeChatSenderType {
	return SenderTypeWxPusher
}

func (s *WxPusherSender) Send(ctx context.Context, content string, summary string, extra map[string]any) error {
	msg := struct {
		AppToken    string   `json:"appToken"`
		Content     string   `json:"content"`
		Summary     string   `json:"summary,omitempty"`
		ContentType int      `json:"contentType"`
		TopicIds    []int64  `json:"topicIds,omitempty"`
		UIds        []string `json:"uids,omitempty"`
	}{
		AppToken:    s.config.AppToken,
		Content:     content,
		Summary:     summary,
		ContentType: 1, // 1=文本
		TopicIds:    s.config.TopicIDs,
	}

	// 如果指定了用户ID，添加到发送目标
	if uid, ok := extra["user_id"].(string); ok {
		msg.UIds = []string{uid}
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message failed: %w", err)
	}

	// 使用配置的 API 地址
	req, err := http.NewRequestWithContext(ctx, "POST", s.config.ApiUrl, bytes.NewReader(body))
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
		Code    int    `json:"code"`
		Msg     string `json:"msg"`
		Success bool   `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response failed: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("send message failed: %s", result.Msg)
	}

	logger.Info("WxPusher message sent successfully",
		zap.Int64s("topic_ids", s.config.TopicIDs))

	return nil
}
