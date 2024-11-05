package sender

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"notify/internal/config"
	"notify/pkg/logger"

	"go.uber.org/zap"
)

type DingTalkSender struct {
	config config.DingTalkConfig
	client *http.Client
}

type DingTalkMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

func NewDingTalkSender(config config.DingTalkConfig) *DingTalkSender {
	return &DingTalkSender{
		config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *DingTalkSender) Send(ctx context.Context, content string, summary string, extra map[string]any) error {
	// 如果有摘要，添加到消息内容前面
	messageContent := content
	if summary != "" {
		messageContent = fmt.Sprintf("【%s】\n\n%s", summary, content)
	}

	msg := DingTalkMessage{
		MsgType: "text",
		Text: struct {
			Content string `json:"content"`
		}{
			Content: messageContent,
		},
	}

	// 生成签名
	timestamp := time.Now().UnixMilli()
	sign := s.generateSign(timestamp)

	// 构建URL
	url := fmt.Sprintf(
		"https://oapi.dingtalk.com/robot/send?access_token=%s&timestamp=%d&sign=%s",
		s.config.AccessToken,
		timestamp,
		sign,
	)

	// 发送消息
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message failed: %w", err)
	}

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

	logger.Info("DingTalk message sent successfully",
		zap.String("content", content))

	return nil
}

func (s *DingTalkSender) generateSign(timestamp int64) string {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, s.config.Secret)
	h := hmac.New(sha256.New, []byte(s.config.Secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
