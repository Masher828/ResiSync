package migrator

import "ResiSync/pkg/models"

func Migrate(migrate models.Migration) {
	migrate.Migrate()
}
