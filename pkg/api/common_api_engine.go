package api

import (
	"ResiSync/pkg/constants"
	postgresclient "ResiSync/pkg/database/postgres"
	redisclient "ResiSync/pkg/database/redis"
	"ResiSync/pkg/logger"
	"ResiSync/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var ApplicationContext *models.ApplicationContextStruct

func init() {
	ApplicationContext = new(models.ApplicationContextStruct)
}

func PrepareApplicationContext() error {

	var err error = nil
	ApplicationContext.Postgres, err = postgresclient.GetPostgresClient()
	if err != nil {
		log.Error("Error while connecting postgres", zap.Error(err))
		return err
	}

	ApplicationContext.Redis, err = redisclient.GetRedisClient()
	if err != nil {
		log.Error("Error while connecting redis", zap.Error(err))
		return err
	}

	return nil
}

func GetAccessToken(c *gin.Context) string {

	authHeader := c.GetHeader("Authorization")

	accessTokenSplit := strings.Split(authHeader, " ")

	accessToken := ""

	if len(accessTokenSplit) != 2 {
		accessToken = accessTokenSplit[1]
	}

	return accessToken
}

func GetUserContextFromAccessToken(requestContext *models.ResiSyncRequestContext, accessToken string) (*models.UserContext, error) {

	redisDB := ApplicationContext.Redis

	log := logger.GetBasicLogger()

	result := redisDB.Get(requestContext.Context, fmt.Sprintf(constants.AccessTokenToUserFormatKey, accessToken))

	if result.Err() != nil {
		log.Error("Error while getting accessToken", zap.String("accessToken", accessToken), zap.Error(result.Err()))
		return nil, result.Err()
	}

	userBytes, err := result.Bytes()
	if err != nil {
		log.Error("Error while getting accessToken", zap.String("accessToken", accessToken), zap.Error(err))
		return nil, err
	}

	var userContext *models.UserContext

	err = json.Unmarshal(userBytes, &userContext)
	if err != nil {
		log.Error("Error while unmarshalling user context", zap.String("access token", accessToken), zap.Error(err))
		return nil, err
	}

	return userContext, nil
}

func SetupRoutes(engine *gin.Engine, routerContext models.RouteContext) {

	routerContext.SetupPrivateRoutes(engine)

	routerContext.SetupPublicRoutes(engine)
}

func GracefulShutdownApp(srv *http.Server, shutdown models.Shutdown) {

	log := logger.GetBasicLogger()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	closeConnections(log)

	if shutdown != nil {
		shutdown.CloseAppSpecificResources()
	}

	stop()

	log.Info("shutting down gracefully, press Ctrl+c again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shut down", zap.Error(err))
	}

	log.Info("Server terminated!!")

}

func closeConnections(log *zap.Logger) {

	postgresclient.ClosePostgresClient(log, ApplicationContext.Postgres)

	redisclient.CloseRedisClient(log, ApplicationContext.Redis)

	if tracer != nil {
		tracer.Shutdown(context.TODO())
	}

	log.Info("Connections closed")
}
