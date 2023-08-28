package controller

import (
	"ResiSync/pkg/api"
	shared_models "ResiSync/shared/models"
	userModels "ResiSync/user/internal/models"
	authfacade "ResiSync/user/internal/services_facade/auth_facade"
	"ResiSync/user/internal/user_errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SignIn(c *gin.Context) {
	requestContext := api.GetRequestContextFromRequest(c)
	span := api.AddTrace(requestContext, "info", "SignIn")
	if span != nil {
		defer span.End()
	}

	log := requestContext.Log

	var user userModels.ResidentDTO
	err := c.ShouldBind(&user)
	if err != nil {
		log.Error("error while binding user signin data", zap.Error(err))
		c.Error(user_errors.ErrInvalidPayload)
		return
	}

	err = authfacade.SignIn(*requestContext, &user)
	if err != nil {
		log.Error("error while signing in the user", zap.Error(err))
		c.Error(err)
		return
	}

	response := shared_models.Response{Status: "ok", Data: user}

	c.JSON(http.StatusOK, response)
}

func SignUp(c *gin.Context) {
	requestContext := api.GetRequestContextFromRequest(c)
	span := api.AddTrace(requestContext, "info", "SignUp")
	if span != nil {
		defer span.End()
	}

	log := requestContext.Log

	var user userModels.ResidentDTO
	err := c.ShouldBind(&user)
	if err != nil {
		log.Error("error while binding user sign up data", zap.Error(err))
		c.Error(user_errors.ErrInvalidPayload)
		return
	}

	err = authfacade.SignUp(*requestContext, &user)
	if err != nil {
		log.Error("error while signing up the user", zap.Error(err))
		c.Status(http.StatusBadRequest)
		c.Error(err)
		return
	}

	response := shared_models.Response{Status: "ok"}

	c.JSON(http.StatusOK, response)
}
