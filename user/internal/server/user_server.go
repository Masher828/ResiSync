package main

import (
	"ResiSync/pkg/api"
	"ResiSync/pkg/logger"
	"ResiSync/pkg/models"
	shared_api "ResiSync/shared/api"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

const (
	host    = "127.0.0.1"
	port    = "8081"
	appName = "User"
)

type Shutdown struct{}

func main() {
	log := logger.GetAppStartupLogger().With(zap.String("appName", appName))

	log.Info("Starting running application")

	applicationContext := models.ApplicationContextStruct{
		AppName: appName,
	}

	engine := api.GetRestApiEngine(&applicationContext)

	engine.Use(shared_api.HandleError())

	srv := &http.Server{
		Addr:    host + ":" + port,
		Handler: engine,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Error("Error while starting server", zap.String("app name", appName), zap.Error(err))
		}
	}()

	api.GracefulShutdownApp(srv, &Shutdown{})

}

func (sd *Shutdown) CloseAppSpecificResources() {
	fmt.Println("Closing app specific resources ", appName)
}
