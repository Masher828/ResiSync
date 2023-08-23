package config

import (
	"ResiSync/pkg/constants"
	pkgerror "ResiSync/pkg/errors"
	"os"
	"strings"

	consul "github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

var consulClient *consul.Client

func getConsulHost() string {
	host := os.Getenv(constants.ConsulHost)

	return strings.TrimSuffix(host, "/")
}

func init() {
	log.Info("Initializing consul client")
	config := consul.DefaultConfig()
	config.Address = getConsulHost()
	config.Token = os.Getenv(constants.ConsulHttpToken)

	var err error = nil
	consulClient, err = consul.NewClient(config)
	if err != nil {
		log.Panic("Error while initializing consul", zap.Error(err))
	}
}

func GetConsulKeyValueList(key string) ([]string, error) {

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
