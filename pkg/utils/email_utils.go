package utils

import (
	pkg_constants "ResiSync/pkg/constants"
	"ResiSync/pkg/logger"
	"ResiSync/pkg/models"
	"ResiSync/pkg/security"
	"html"
	"net/smtp"

	"github.com/go-playground/validator/v10"
	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func getSmtpConfig() (*models.SmtpConfig, error) {
	log := logger.GetBasicLogger()

	key := pkg_constants.ConfigSmtpKey

	var smtpConfig models.SmtpConfig

	err := viper.UnmarshalKey(key, &smtpConfig)
	if err != nil {
		log.Error("Error while unmarshalling smtp config", zap.Error(err))
		return nil, err
	}

	validate := validator.New()

	err = validate.Struct(smtpConfig)
	if err != nil {
		log.Error("Error while validating smtp config", zap.Error(err))
		return nil, err
	}

	return &smtpConfig, nil
}

func getSmtpDetails() (smtp.Auth, string, error) {
	log := logger.GetBasicLogger()

	smtpconfig, err := getSmtpConfig()
	if err != nil {
		log.Error("Error while fetching smtp config", zap.Error(err))
		return nil, "", err
	}

	decryptedPassword, err := security.DecryptPassword(smtpconfig.Password, smtpconfig.PasswordNonce)
	if err != nil {
		log.Error("Error while decryption smtp password", zap.Error(err))
		return nil, "", err
	}

	auth := smtp.PlainAuth("", smtpconfig.Username, decryptedPassword, smtpconfig.Host)

	return auth, smtpconfig.Host + ":" + smtpconfig.Port, nil
}

func SendEmail(from, to, subject, body string) error {
	return SendEmailToMultiple([]string{to}, from, subject, body)
}

func SendEmailToMultiple(to []string, from, subject, body string) error {
	log := logger.GetBasicLogger()
	auth, host, err := getSmtpDetails()
	if err != nil {
		log.Error("Error while getting smtp details", zap.Error(err))
		return err
	}

	e := email.NewEmail()
	e.From = from
	e.To = to
	e.Subject = subject
	e.HTML = []byte(html.UnescapeString(body))

	err = e.Send(host, auth)
	if err != nil {
		log.Error("Error while sending email", zap.Error(err))
		return err
	}

	return nil
}
