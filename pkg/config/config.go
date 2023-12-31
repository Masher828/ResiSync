package config

import (
	pkg_constants "ResiSync/pkg/constants"
	"ResiSync/pkg/logger"
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var log *zap.Logger

func init() {
	log = logger.GetAppStartupLogger()
}

func LoadEnv() error {
	viper.SetConfigFile(os.Getenv(pkg_constants.EnvFilePathKey))

	return viper.ReadInConfig()
}

func LoadConfig() error {
	log.Info("started loading config")

	viper.SetConfigType("yaml")

	key := strings.Join([]string{viper.GetString(pkg_constants.ConsulConfigKey), pkg_constants.CommonConfigFolderName}, "/")

	keyValueList, err := GetConsulKeyValueList(key)
	if err != nil {
		log.Error("Error while fetching from consul", zap.String("key", key), zap.Error(err))
		return err
	}

	for _, value := range keyValueList {
		err = viper.MergeConfig(strings.NewReader(value))
		if err != nil {
			log.Error("Error in populating viper from consul", zap.Error(err))
			return err
		}
	}

	return nil
}
