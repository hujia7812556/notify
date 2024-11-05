package parser

import (
	"encoding/json"
	"errors"
)

type Platform string

const (
	PlatformWeChat   Platform = "wechat"
	PlatformDingTalk Platform = "dingtalk"
)

type Message struct {
	Platform Platform       `json:"platform"`
	Content  string         `json:"content"`
	Summary  string         `json:"summary,omitempty"`
	Extra    map[string]any `json:"extra,omitempty"`
}

func Parse(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	if err := msg.Validate(); err != nil {
		return nil, err
	}

	return &msg, nil
}

// Validate 验证消息格式
func (m *Message) Validate() error {
	if m.Platform == "" {
		return errors.New("platform is required")
	}
	if m.Content == "" {
		return errors.New("content is required")
	}
	if !isValidPlatform(m.Platform) {
		return errors.New("unsupported platform")
	}
	return nil
}

// isValidPlatform 检查平台是否支持
func isValidPlatform(platform Platform) bool {
	switch platform {
	case PlatformWeChat, PlatformDingTalk:
		return true
	default:
		return false
	}
}
