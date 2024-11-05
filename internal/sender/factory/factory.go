package factory

import (
	"fmt"
	"notify/internal/config"
	"notify/internal/sender"
	"notify/internal/sender/wechat"
)

// CreateWeChatSender 创建微信发送器
func CreateWeChatSender(senderType wechat.WeChatSenderType, config config.WeChatConfig) (sender.Sender, error) {
	switch senderType {
	case wechat.SenderTypeWeCom:
		return wechat.NewWeComSender(config.WeCom), nil
	case wechat.SenderTypeWxPusher:
		return wechat.NewWxPusherSender(config.WxPusher), nil
	default:
		return nil, fmt.Errorf("unsupported WeChat sender type: %s", senderType)
	}
}
