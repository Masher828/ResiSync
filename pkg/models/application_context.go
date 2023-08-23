package models

import (
	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5"
)

type ApplicationContext struct {
	AppName  string
	Postgres *pgx.Conn
	Redis    *redis.Client
}

type LoggerConfig struct {
	Level string `mapstructure:"level" validate:"required, oneof=debug info warn error panic fatal"`
	Env   string `mapstructure:"env" validate:"required,oneof=production development"`
}

type SmtpConfig struct {
	Host          string `mapstructure:"host" validate:"required"`
	Username      string `mapstructure:"username" validate:"required"`
	Password      string `mapstructure:"password" validate:"required"`
	PasswordKey   string `mapstructure:"password_key" validate:"required"`
	PasswordNonce string `mapstructure:"password_nonce" validate:"required"`
	Port          string `mapstructure:"port" validate:"required"`
}
