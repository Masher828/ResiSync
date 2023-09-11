package controller

import (
	"ResiSync/app/internal/models"
	user_models "ResiSync/app/internal/models/user"
	"ResiSync/app/internal/services/auth_service"
	auth_facade "ResiSync/app/internal/services_facade/auth_facade"
	"ResiSync/pkg/api"
	shared_errors "ResiSync/shared/errors"
	shared_models "ResiSync/shared/models"
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

	var user user_models.ResidentDTO
	err := c.ShouldBind(&user)
	if err != nil {
		log.Error("error while binding user signin data", zap.Error(err))
		c.Error(shared_errors.ErrInvalidPayload)
		return
	}

	err = auth_facade.SignIn(*requestContext, &user)
	if err != nil {
		log.Error("error while signing in the user", zap.Error(err))
		c.Error(err)
		return
	}

	response := shared_models.Response{Status: "ok", Data: user, StatusCode: http.StatusOK}

	c.JSON(http.StatusOK, response)
}

func SignUp(c *gin.Context) {
	requestContext := api.GetRequestContextFromRequest(c)
	span := api.AddTrace(requestContext, "info", "SignUp")
	if span != nil {
		defer span.End()
	}

	log := requestContext.Log

	var user user_models.ResidentDTO
	err := c.ShouldBind(&user)
	if err != nil {
		log.Error("error while binding user sign up data", zap.Error(err))
		c.Error(shared_errors.ErrInvalidPayload)
		return
	}

	err = auth_facade.SignUp(*requestContext, &user)
	if err != nil {
		log.Error("error while signing up the user", zap.Error(err))
		c.Status(http.StatusBadRequest)
		c.Error(err)
		return
	}

	response := shared_models.Response{Status: "ok", StatusCode: http.StatusOK}

	c.JSON(http.StatusOK, response)
}

func LogOut(c *gin.Context) {
	requestContext := api.GetRequestContextFromRequest(c)
	span := api.AddTrace(requestContext, "info", "LogOut")
	if span != nil {
		defer span.End()
	}

	go auth_service.LogOut(*requestContext)
	response := shared_models.Response{Status: "ok", StatusCode: http.StatusOK}

	c.JSON(http.StatusOK, response)
}

func ResetPassword(c *gin.Context) {
	requestContext := api.GetRequestContextFromRequest(c)
	span := api.AddTrace(requestContext, "info", "ResetPassword")
	defer span.End()

	log := requestContext.Log

	var resetPassword models.ResetPassword
	if err := c.Bind(&resetPassword); err != nil {
		log.Error("error while unmarshalling request body", zap.Error(err))
		c.Error(shared_errors.ErrInvalidPayload)
		return
	}

	err := auth_facade.ResetPassword(*requestContext, &resetPassword)
	if err != nil {
		log.Error("error while resetting password", zap.Error(err))
		c.Error(err)
		return
	}

	response := shared_models.Response{Status: "ok", StatusCode: http.StatusOK}

	c.JSON(http.StatusOK, response)

}
