package types

// LogConfig 日志配置结构体
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
	Debug  struct {
		Output string `mapstructure:"output"` // debug模式下的输出
	} `mapstructure:"debug"`
	Release struct {
		Output string `mapstructure:"output"` // release模式下的输出
	} `mapstructure:"release"`
}
