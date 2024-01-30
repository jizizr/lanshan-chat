package config

type AuthConfig struct {
	JwtConfig *JwtConfig `mapstructure:"jwt" json:"jwt"`
}

type JwtConfig struct {
	SecretKey   string `mapstructure:"secret-key" json:"secret-key"`
	ExpiresTime int64  `mapstructure:"expires-time" json:"expires-time"`
}
