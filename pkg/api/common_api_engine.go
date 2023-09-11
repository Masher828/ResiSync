package api

import (
	aws_services "ResiSync/pkg/aws"
	pkg_constants "ResiSync/pkg/constants"
	postgresclient "ResiSync/pkg/database/postgres"
	redisclient "ResiSync/pkg/database/redis"
	"ResiSync/pkg/logger"
	pkg_models "ResiSync/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	openotel "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var ApplicationContext *pkg_models.ApplicationContextStruct

func init() {
	ApplicationContext = new(pkg_models.ApplicationContextStruct)
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

	ApplicationContext.S3Session, err = aws_services.CreateNewS3Session()
	if err != nil {
		log.Error("Error while creating s3 session", zap.Error(err))
		return err
	}
	return nil
}

func GetAccessToken(c *gin.Context) string {

	authHeader := c.GetHeader("Authorization")

	accessTokenSplit := strings.Split(authHeader, " ")

	accessToken := ""

	if len(accessTokenSplit) >= 2 {
		accessToken = accessTokenSplit[1]
	}

	return accessToken
}

func GetUserContextFromAccessToken(requestContext *pkg_models.ResiSyncRequestContext, accessToken string) (*pkg_models.UserContext, error) {

	redisDB := ApplicationContext.Redis

	log := logger.GetBasicLogger()

	result := redisDB.Get(requestContext.Context, fmt.Sprintf(pkg_constants.AccessTokenToUserFormatKey, accessToken))

	if result.Err() != nil {
		log.Error("Error while getting accessToken", zap.String("accessToken", accessToken), zap.Error(result.Err()))
		return nil, result.Err()
	}

	userBytes, err := result.Bytes()
	if err != nil {
		log.Error("Error while getting accessToken", zap.String("accessToken", accessToken), zap.Error(err))
		return nil, err
	}

	var userContext *pkg_models.UserContext

	err = json.Unmarshal(userBytes, &userContext)
	if err != nil {
		log.Error("Error while unmarshalling user context", zap.String("access token", accessToken), zap.Error(err))
		return nil, err
	}

	return userContext, nil
}

func SetupRoutes(engine *gin.Engine, routerContext pkg_models.RouteContext) {

	routerContext.SetupPrivateRoutes(engine)

	routerContext.SetupPublicRoutes(engine)
}

func GracefulShutdownApp(srv *http.Server, shutdown pkg_models.Shutdown) {

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

func GetRequestContextFromRequest(c *gin.Context) *pkg_models.ResiSyncRequestContext {
	requestContext, _ := c.Get(pkg_constants.RequestContextKey)

	return requestContext.(*pkg_models.ResiSyncRequestContext)
}

func AddTrace(requestContext *pkg_models.ResiSyncRequestContext, level, spanName string) trace.Span {

	tracerContext, span := openotel.Tracer("").Start(requestContext.Context, spanName)

	span.SetAttributes(attribute.String("traceID", requestContext.GetTraceID()))

	if requestContext.GetUserContext() != nil {
		span.SetAttributes(attribute.Int64("userContext", requestContext.GetUserContext().ID))
	}

	span.SetAttributes(attribute.String("routePath", requestContext.GetRoutePath()))
	requestContext.Context = tracerContext
	return span

}
