package userMigrator

import (
	"ResiSync/pkg/logger"
)

type Migrator struct{}

func (migrator *Migrator) Migrate() {
	// basic structure will be created via docker
	// after that everything will come here
	// db := api.ApplicationContext.Postgres

	log := logger.GetBasicLogger()

	log.Info("Started migration")
	log.Info("Complted migration")

}
