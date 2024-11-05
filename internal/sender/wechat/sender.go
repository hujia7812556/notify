package wechat

import (
	"context"
)

// WeChatSenderType 定义微信发送器类型
type WeChatSenderType string

const (
	SenderTypeWeCom    WeChatSenderType = "wecom"    // 企业微信
	SenderTypeWxPusher WeChatSenderType = "wxpusher" // WxPusher
)

// WeChatSender 微信发送器接口
type WeChatSender interface {
	Type() WeChatSenderType
	Send(ctx context.Context, content string, extra map[string]any) error
}
