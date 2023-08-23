package main

import (
	"ResiSync/pkg/api"
	"ResiSync/pkg/logger"
	"ResiSync/pkg/models"

	"go.uber.org/zap"
)

const (
	host    = "127.0.0.1"
	port    = "8081"
	appName = "User"
)

func main() {
	log := logger.GetAppStartupLogger().With(zap.String("appName", appName))

	log.Info("Starting running application")

	applicationContext := models.ApplicationContext{
		AppName: appName,
	}

	engine := api.GetRestApiEngine(&applicationContext)

	engine.Run(":8888")

}
