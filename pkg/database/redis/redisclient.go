package redisclient

import (
	"ResiSync/pkg/constants"
	"ResiSync/pkg/models"
	"ResiSync/pkg/otel"
	"ResiSync/pkg/security"
	"context"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

var tracer *sdktrace.TracerProvider = nil

func GetRedisClient() (*redis.Client, error) {

	var redisConf models.Redis
	err := viper.UnmarshalKey(constants.ConfigSectionRedis, &redisConf)
	if err != nil {
		log.Println("Error while unmarshalling redis config", err)
		return nil, err
	}

	validate := validator.New()

	err = validate.Struct(redisConf)
	if err != nil {
		log.Println("Error while validating redis conf", err)
		return nil, err
	}

	password, err := security.DecryptPassword(redisConf.Password, redisConf.PasswordNonce)
	if err != nil {
		log.Println("Error while decrypting redis password", err)
		return nil, err
	}
	opt := redis.Options{
		Addr:     redisConf.Address,
		Password: password,
		DB:       redisConf.Db,
	}

	client := redis.NewClient(&opt)

	err = client.Ping(context.Background()).Err()
	if err != nil {
		log.Println("Error while pinging redis", err)
		return nil, err
	}

	if otel.IsTracingEnabled() {
		tracer, err = otel.GetTracerProvider("redis")
		if err != nil {
			log.Println("Error while getting tracer provider for redis", err)
			return nil, err
		}

		err = redisotel.InstrumentTracing(client, redisotel.WithTracerProvider(tracer))
		if err != nil {
			log.Println("Error while instrumenting redis", err)
			return nil, err
		}

	}
	return client, nil

}

func closeTracer() {
	if tracer != nil {
		tracer.Shutdown(context.TODO())
	}
}

func CloseRedisClient(log *zap.Logger, conn *redis.Client) {

	closeTracer()

	if conn == nil {
		log.Info("redis Connection is nil. Returning")
		return
	}

	err := conn.Close()
	if err != nil {
		log.Error("Error while closing redis object ", zap.Error(err))
		return
	}
}
