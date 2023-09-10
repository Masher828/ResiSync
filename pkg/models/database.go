package models

type Postgres struct {
	Host          string `mapstructure:"host" validate:"required"`
	Port          string `mapstructure:"port" validate:"required"`
	Username      string `mapstructure:"username" validate:"required"`
	Password      string `mapstructure:"password" validate:"required,base64"`
	PasswordNonce string `mapstructure:"password_nonce" validate:"required,base64"`
	Db            string `mapstructure:"db" validate:"required"`
}

type Redis struct {
	Address       string `mapstructure:"address" validate:"required"`
	Password      string `mapstructure:"password" validate:"required,base64"`
	PasswordNonce string `mapstructure:"password_nonce" validate:"required,base64"`
	Db            int    `mapstructure:"db"`
}
