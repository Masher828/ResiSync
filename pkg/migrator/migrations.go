package migrator

import pkg_models "ResiSync/pkg/models"

func Migrate(migrate pkg_models.Migration) {
	migrate.Migrate()
}
