package config

import (
	"notify/pkg/logger"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig
	Dispatcher  DispatcherConfig
	WeChat      WeChatConfig
	DingTalk    DingTalkConfig
	Log         logger.LogConfig
	HealthCheck HealthCheckConfig
}

type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	Token        string        `mapstructure:"token"`
}

type DispatcherConfig struct {
	BufferSize     int `mapstructure:"buffer_size"`
	WorkerPoolSize int `mapstructure:"worker_pool_size"`
}

type WeChatConfig struct {
	SenderType string         `mapstructure:"sender_type"`
	WeCom      WeComConfig    `mapstructure:"wecom"`
	WxPusher   WxPusherConfig `mapstructure:"wxpusher"`
}

type WeComConfig struct {
	CorpID    string `mapstructure:"corp_id"`
	AgentID   string `mapstructure:"agent_id"`
	AppSecret string `mapstructure:"app_secret"`
}

type WxPusherConfig struct {
	AppToken string  `mapstructure:"app_token"`
	TopicIDs []int64 `mapstructure:"topic_ids"`
	QPS      int     `mapstructure:"qps"`
	ApiUrl   string  `mapstructure:"api_url"`
}

type DingTalkConfig struct {
	AccessToken string `mapstructure:"access_token"`
	Secret      string `mapstructure:"secret"`
}

type HealthCheckConfig struct {
	Enabled   bool          `mapstructure:"enabled"`
	CheckTime string        `mapstructure:"check_time"`
	Timeout   time.Duration `mapstructure:"timeout"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 按优先级顺序添加配置路径
	viper.AddConfigPath("/etc/notify")   // 系统配置目录
	viper.AddConfigPath("$HOME/.notify") // 用户目录
	viper.AddConfigPath("./etc")         // 当前目录的etc子目录
	viper.AddConfigPath("./config")      // 当前目录的config子目录

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
