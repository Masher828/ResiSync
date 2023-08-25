package models

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type key string

const (
	requestContextTraceIDKey     key = "trace_id"
	requestContextUserContextKey key = "user_context"
)

type UserContext struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	AccessToken string `json:"accessToken"`
}

type ResiSyncRequestContext struct {
	Context context.Context
	Log     *zap.Logger
}

func (rs *ResiSyncRequestContext) SetTraceID(traceID string) {
	rs.Context = context.WithValue(rs.Context, requestContextTraceIDKey, traceID)
}

func (rs *ResiSyncRequestContext) GetTraceID() string {
	if val, ok := rs.Context.Value(requestContextTraceIDKey).(string); ok {
		return val
	}
	return ""
}

func (rs *ResiSyncRequestContext) SetUserContext(userContext *UserContext) {
	rs.Context = context.WithValue(rs.Context, requestContextTraceIDKey, userContext)
}

func (rs *ResiSyncRequestContext) GetUserContext() *UserContext {
	if val, ok := rs.Context.Value(requestContextTraceIDKey).(*UserContext); ok {
		return val
	}
	return nil
}

type RouteContext interface {
	SetupPublicRoutes(*gin.Engine)

	SetupPrivateRoutes(*gin.Engine)
}

type Shutdown interface {
	CloseAppSpecificResources()
}
