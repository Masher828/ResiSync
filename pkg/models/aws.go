package models

type Aws struct {
	AccessKeyId        string `mapstructure:"access_key_id" validate:"required"`
	EncryptedSecretKey string `mapstructure:"encrypted_secret_key" validate:"required,base64"`
	SecretKeyNonce     string `mapstructure:"secret_key_nonce" validate:"required,base64"`
	Token              string `mapstructure:"token"`
	Region             string `mapstructure:"region" validate:"required"`
}
