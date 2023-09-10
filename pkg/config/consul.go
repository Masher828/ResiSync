package config

import (
	pkg_constants "ResiSync/pkg/constants"
	pkgerror "ResiSync/pkg/errors"
	"strings"

	consul "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var consulClient *consul.Client = nil

func getConsulHost() string {
	host := viper.GetString(pkg_constants.ConsulHost)
	return strings.TrimSuffix(host, "/")
}

func getConsulClient() *consul.Client {
	if consulClient == nil {
		log.Info("Initializing consul client")
		config := consul.DefaultConfig()
		config.Address = getConsulHost()
		config.Token = viper.GetString(pkg_constants.ConsulHttpToken)

		var err error = nil
		consulClient, err = consul.NewClient(config)
		if err != nil {
			log.Panic("Error while initializing consul", zap.Error(err))
		}
	}

	return consulClient
}

func GetConsulKeyValueList(key string) ([]string, error) {

	consulClient = getConsulClient()

	kvPair, _, err := consulClient.KV().List(key, nil)
	if err != nil {
		log.Error("Error while fetching from consul", zap.String("key", key), zap.Error(err))
		return nil, err
	}

	if kvPair == nil {
		err = pkgerror.ErrKeyDoesNotExist
		log.Error("", zap.Error(err))
		return nil, err
	}

	var values []string
	for _, value := range kvPair {
		values = append(values, string(value.Value))
	}

	return values, nil
}
