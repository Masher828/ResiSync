package postgres_db_test

import (
	"ResiSync/pkg/api"
	"ResiSync/pkg/config"
	"testing"
)

func init() {
	config.LoadEnv()

	config.LoadConfig()
}

func Test_Postgres(t *testing.T) {
	db := api.ApplicationContext.Postgres
	var id int64
	err := db.Exec("SELECT nextval('hibernate_sequence');").Scan(&id).Error

	t.Error(id, err)
}
