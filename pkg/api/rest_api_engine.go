package api

import (
	"ResiSync/pkg/config"
	pkg_constants "ResiSync/pkg/constants"
	"ResiSync/pkg/logger"
	pkg_models "ResiSync/pkg/models"
	"ResiSync/pkg/otel"
	pkg_utils "ResiSync/pkg/utils"
	"context"
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

func GetRestApiEngine(appContext *pkg_models.ApplicationContextStruct) *gin.Engine {

	err := config.LoadEnv()
	if err != nil {
		log.Panic("Error while loading env file", zap.String("appName", appContext.AppName), zap.Error(err))
		return nil
	}

	err = config.LoadConfig()
	if err != nil {
		log.Panic("Error while loading config", zap.String("appName", appContext.AppName), zap.Error(err))
		return nil
	}

	logger.InitializeLoggerWithHook(pkg_utils.GetLoggerEmailHook(appContext.AppName))

	log := logger.GetBasicLogger()

	log.Info("Logger intialized")

	if otel.IsTracingEnabled() {
		tracer, err = otel.InitTracer(appContext.AppName)
		if err != nil {
			log.Error("error while enabling tracer", zap.Error(err), zap.String("appname", appContext.AppName))
		} else if tracer == nil {
			log.Info("Tracer is not enabled")
		}
	}

	//TODO add metrics

	if os.Getenv(pkg_constants.EnvEnvironment) == pkg_constants.EnvironmentProduction {
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
			requestContextInterface, _ := c.Get(pkg_constants.RequestContextKey)
			requestContext := requestContextInterface.(*pkg_models.ResiSyncRequestContext)

			userContext, err := GetUserContextFromAccessToken(requestContext, accessToken)
			if err != nil {

				log.Error("Error while getting usercontext", zap.String("access token", accessToken), zap.Error(err))
				c.Set(pkg_constants.RequestAuthenticatedKey, false)

			} else {

				requestContext.SetUserContext(userContext)
				// requestContext.Log = requestContext.Log.With(zap.Field{Key: "user Id", Integer: userContext.ID})
				c.Set(pkg_constants.RequestUserContextKey, userContext)
				c.Set(pkg_constants.RequestAuthenticatedKey, true)
				c.Set(pkg_constants.RequestContextKey, requestContext)
			}

		} else {
			c.Set(pkg_constants.RequestAuthenticatedKey, false)
		}

		c.Next()
	}
}

func addRequestContext(appName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// requestContextBytes := c.Request.Header.Get(constants.RequestContextHeaderKey)

		// TODO add request context for calling service internally

		var requestContext pkg_models.ResiSyncRequestContext

		requestContext.Context = context.TODO()

		requestContext.Log = logger.GetBasicLogger()

		requestContext.SetTraceID(uuid.New().String())

		c.Set(pkg_constants.RequestContextKey, &requestContext)

		c.Next()

	}
}
