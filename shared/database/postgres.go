package postgres_db

import (
	"ResiSync/pkg/api"
	"ResiSync/pkg/models"

	"go.uber.org/zap"
)

func GetSequenceId(requestContext models.ResiSyncRequestContext) (int64, error) {
	span := api.AddTrace(&requestContext, "info", "GetSequenceId")
	defer span.End()

	log := requestContext.Log

	db := api.ApplicationContext.Postgres

	var id int64
	err := db.Raw(`SELECT nextval($1);`, "hibernate_sequence").Scan(&id).Error
	if err != nil {
		log.Error("Error while getting hibernate_sequence", zap.Error(err))
		return 0, err
	}

	return id, nil
}

func SaveOrUpdate(requestContext models.ResiSyncRequestContext, data interface{}) error {
	span := api.AddTrace(&requestContext, "info", "SaveOrUpdate")
	defer span.End()

	log := requestContext.Log

	postgresDB := api.ApplicationContext.Postgres

	err := postgresDB.Save(data).Error
	if err != nil {
		log.Error("error while updating data", zap.Error(err))
		return err
	}
	return nil
}

func GetWithFields(requestContext models.ResiSyncRequestContext, data interface{}, fields ...string) error {
	span := api.AddTrace(&requestContext, "info", "GetWithFields")
	defer span.End()

	log := requestContext.Log

	postgresDB := api.ApplicationContext.Postgres

	err := postgresDB.Select(fields).Find(data).Error
	if err != nil {
		log.Error("error while getting data", zap.Error(err))
		return err
	}
	return nil

}

func GetData(requestContext models.ResiSyncRequestContext, data interface{}) error {
	span := api.AddTrace(&requestContext, "info", "GetData")
	defer span.End()

	log := requestContext.Log

	postgresDB := api.ApplicationContext.Postgres

	err := postgresDB.Find(data).Error
	if err != nil {
		log.Error("error while getting data", zap.Error(err))
		return err
	}
	return nil

}

func UpdateWithFields(requestContext models.ResiSyncRequestContext, data interface{}, fields ...string) error {
	span := api.AddTrace(&requestContext, "info", "UpdateWithFields")
	defer span.End()

	log := requestContext.Log

	postgresDB := api.ApplicationContext.Postgres

	err := postgresDB.Model(data).Select(fields).Updates(data).Error
	if err != nil {
		log.Error("error while updating data", zap.Error(err))
		return err
	}
	return nil
}
