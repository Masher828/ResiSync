package postgresclient

import (
	"ResiSync/pkg/constants"
	"ResiSync/pkg/models"
	"ResiSync/pkg/otel"
	"ResiSync/pkg/security"
	"context"
	"fmt"
	"log"

	"github.com/spf13/viper"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

var tracer *sdktrace.TracerProvider

func GetPostgresClient() (*gorm.DB, error) {

	var postgresConf models.Postgres

	err := viper.UnmarshalKey(constants.ConfigSectionPostgres, &postgresConf)
	if err != nil {
		log.Panic("Error while unmarshalling postgres config", err)
	}

	password, err := security.DecryptPassword(postgresConf.Password, postgresConf.PasswordNonce)
	if err != nil {
		log.Panic("Error while decrypting postgres password", err)
	}

	databaseUrl := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", postgresConf.Host, postgresConf.Username, password, postgresConf.Db, postgresConf.Port)

	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
	if err != nil {
		log.Panic("Error while parsing postgres databaseurl", err)
	}

	tracer, err = otel.GetTracerProvider("postgres")
	if err != nil {
		log.Panic("Error while parsing postgres databaseurl", err)
	}

	db.Use(tracing.NewPlugin(tracing.WithTracerProvider(tracer)))

	return db, nil
}

func closeTracer() {
	if tracer != nil {
		tracer.Shutdown(context.TODO())
	}
}

func ClosePostgresClient(log *zap.Logger, conn *gorm.DB) {
	closeTracer()
	if conn == nil {
		log.Info("Gorm Connection is nil. Returning")
		return
	}

	db, err := conn.DB()
	if err != nil {
		log.Error("Error while getting DB object from gorm", zap.Error(err))
		return
	}

	err = db.Close()
	if err != nil {
		log.Error("Error while closing DB object from gorm", zap.Error(err))
		return
	}
}
