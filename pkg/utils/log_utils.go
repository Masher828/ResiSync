package utils

import (
	pkg_constants "ResiSync/pkg/constants"
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

func GetLoggerEmailHook(appname string) func(entry zapcore.Entry) error {
	return func(entry zapcore.Entry) error {
		if entry.Level > zapcore.ErrorLevel {
			to := viper.GetString(pkg_constants.SendErrorEmailToKey)
			from := viper.GetString(pkg_constants.SendErrorEmailFromKey)

			if viper.GetBool(pkg_constants.SendErrorEmailKey) && len(to) > 0 && len(from) > 0 {
				env := os.Getenv(pkg_constants.EnvEnvironment)
				body := entry.Message + "\n\n" + entry.Stack
				subject := strings.ToUpper(env) + "-" + appname + "Error"
				go SendEmail(from, to, subject, body)
			}
		}
		return nil
	}
}
