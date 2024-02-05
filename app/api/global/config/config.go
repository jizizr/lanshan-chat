package config

type Config struct {
	ZapConfig      ZapConfig      `mapstructure:"zap_config" json:"zap_config"`
	DatabaseConfig DatabaseConfig `mapstructure:"database" json:"database"`
	ServerConfig   ServerConfig   `mapstructure:"server" json:"server"`
	AuthConfig     AuthConfig     `mapstructure:"auth" json:"auth"`
	FilterConfig   SenFilter      `mapstructure:"filter" json:"filter"`
}
