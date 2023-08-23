package api

import (
	"ResiSync/pkg/config"
	"ResiSync/pkg/logger"
	"ResiSync/pkg/models"
	"ResiSync/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var log *zap.Logger

func init() {
	log = logger.GetAppStartupLogger()
}

func GetRestApiEngine(appContext *models.ApplicationContext) *gin.Engine {

	err := config.LoadConfig()
	if err != nil {
		log.Panic("Error while loading config", zap.String("appName", appContext.AppName), zap.Error(err))
		return nil
	}

	logger.InitializeLoggerWithHook(utils.GetLoggerEmailHook(appContext.AppName))

	log := logger.GetBasicLogger()

	log.Info("Logger intialized")

	return gin.New()

}
