package controller

import (
	user_errors "ResiSync/app/internal/app_errors"
	"ResiSync/app/internal/services/user_service.go"
	"ResiSync/pkg/api"
	shared_models "ResiSync/shared/models"
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func GetUserProfile(c *gin.Context) {
	requestContext := api.GetRequestContextFromRequest(c)
	span := api.AddTrace(requestContext, "info", "LogOut")
	if span != nil {
		defer span.End()
	}

	log := requestContext.Log
	user, err := user_service.GetUserProfile(*requestContext)
	if err != nil {
		log.Error("error while fetching user", zap.Error(err))
		c.Error(user_errors.ErrInvalidPayload)
		return
	}

	response := shared_models.Response{Status: "ok", Data: user}

	c.JSON(http.StatusOK, response)
}
