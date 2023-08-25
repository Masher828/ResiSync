package api

import (
	"ResiSync/pkg/config"
	"ResiSync/pkg/constants"
	"ResiSync/pkg/logger"
	"ResiSync/pkg/models"
	"ResiSync/pkg/otel"
	"ResiSync/pkg/utils"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

var log *zap.Logger

var tracer *sdktrace.TracerProvider

func init() {
	log = logger.GetAppStartupLogger()
}

func GetRestApiEngine(appContext *models.ApplicationContextStruct) *gin.Engine {

	err := config.LoadConfig()
	if err != nil {
		log.Panic("Error while loading config", zap.String("appName", appContext.AppName), zap.Error(err))
		return nil
	}

	logger.InitializeLoggerWithHook(utils.GetLoggerEmailHook(appContext.AppName))

	log := logger.GetBasicLogger()

	log.Info("Logger intialized")

	tracer, err = otel.InitTracer(appContext.AppName)
	if err != nil {
		log.Error("error while enabling tracer", zap.Error(err), zap.String("appname", appContext.AppName))
	} else if tracer == nil {
		log.Info("Tracer is not enabled")
	}

	//TODO add metrics

	if os.Getenv(constants.EnvEnvironment) == constants.EnvironmentProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	addMiddlewares(router, appContext.AppName)

	err = PrepareApplicationContext()
	if err != nil {
		log.Panic("Error while preparing application context", zap.String("appName", appContext.AppName), zap.Error(err))
		return nil
	}

	return router

}

func addMiddlewares(router *gin.Engine, appName string) {

	router.Use(gin.Recovery())

	config := cors.DefaultConfig()

	config.AllowOrigins = []string{"*"}

	router.Use(cors.New(config))

	router.Use(addRequestContext(appName))

	router.Use(applyAuth())

}

func applyAuth() gin.HandlerFunc {

	return func(c *gin.Context) {
		log := logger.GetBasicLogger()

		//implement auth / rbac with loading urls into redis

		accessToken := GetAccessToken(c)

		if len(accessToken) > 0 {
			requestContextInterface, _ := c.Get(constants.RequestContextKey)
			requestContext := requestContextInterface.(*models.ResiSyncRequestContext)

			userContext, err := GetUserContextFromAccessToken(requestContext, accessToken)
			if err != nil {

				log.Error("Error while getting usercontext", zap.String("access token", accessToken), zap.Error(err))
				c.Set(constants.RequestAuthenticatedKey, false)

			} else {

				requestContext.SetUserContext(userContext)
				requestContext.Log = requestContext.Log.With(zap.Field{Key: "user Id", Integer: userContext.ID})
				c.Set(constants.RequestUserContextKey, userContext)
				c.Set(constants.RequestAuthenticatedKey, true)
				c.Set(constants.RequestContextKey, requestContext)
			}

		} else {
			c.Set(constants.RequestAuthenticatedKey, false)
		}

		c.Next()
	}
}

func addRequestContext(appName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// requestContextBytes := c.Request.Header.Get(constants.RequestContextHeaderKey)

		// TODO add request context for calling service internally

		var requestContext models.ResiSyncRequestContext

		requestContext.Log = logger.GetBasicLogger().With(
			zap.Field{Key: "url", String: c.Request.RequestURI})

		requestContext.SetTraceID(uuid.New().String())

		c.Set(constants.RequestContextKey, requestContext)

		c.Next()

	}
}
