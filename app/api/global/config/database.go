package config

import (
	"fmt"
	"time"
)

type DatabaseConfig struct {
	MysqlConfig `mapstructure:"mysql" json:"mysql"`
	RedisConfig `mapstructure:"redis" json:"redis"`
}

type MysqlConfig struct {
	Addr            string        `mapstructure:"addr" json:"addr"`
	Port            string        `mapstructure:"port" json:"port"`
	DB              string        `mapstructure:"db" json:"db"`
	Username        string        `mapstructure:"username" json:"username"`
	Password        string        `mapstructure:"password" json:"password"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" json:"conn_max_idle_time"`
	MaxIdleConn     int           `mapstructure:"max_idle_conns" json:"max_idle_conns"`
	MaxOpenConn     int           `mapstructure:"max_open_conns" json:"max_open_conns"`
	Charset         string        `mapstructure:"charset" json:"charset"`
	Place           string        `mapstructure:"place" json:"place"`
}

type RedisConfig struct {
	Host     string `mapstructure:"addr" json:"addr"`
	Port     string `mapstructure:"port" json:"port"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	DB       int    `mapstructure:"db" json:"db"`
}

func (c *MysqlConfig) GetDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=%s",
		c.Username,
		c.Password,
		c.Addr,
		c.Port,
		c.DB,
		c.Charset,
		c.Place,
	)
}
