package pkg_models

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ApplicationContextStruct struct {
	AppName   string
	Postgres  *gorm.DB
	Redis     *redis.Client
	S3Session *s3.S3
}

type LoggerConfig struct {
	Level string `mapstructure:"level" validate:"required,oneof=debug info warn error panic fatal"`
	Env   string `mapstructure:"env" validate:"required,oneof=production development"`
}

type SmtpConfig struct {
	Host          string `mapstructure:"host" validate:"required"`
	Username      string `mapstructure:"username" validate:"required"`
	Password      string `mapstructure:"password" validate:"required"`
	PasswordNonce string `mapstructure:"password_nonce" validate:"required"`
	Port          string `mapstructure:"port" validate:"required"`
}
