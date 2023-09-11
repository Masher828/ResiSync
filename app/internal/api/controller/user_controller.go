package controller

import (
	user_models "ResiSync/app/internal/models/user"
	"ResiSync/app/internal/services/user_service.go"
	"ResiSync/app/internal/services_facade/user_facade"
	"ResiSync/pkg/api"
	shared_errors "ResiSync/shared/errors"
	shared_models "ResiSync/shared/models"
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func GetUserProfile(c *gin.Context) {
	requestContext := api.GetRequestContextFromRequest(c)
	span := api.AddTrace(requestContext, "info", "GetUserProfile")
	if span != nil {
		defer span.End()
	}

	log := requestContext.Log
	user, err := user_service.GetUserProfile(*requestContext)
	if err != nil {
		log.Error("error while fetching user", zap.Error(err))
		c.Error(shared_errors.ErrInvalidPayload)
		return
	}

	response := shared_models.Response{Status: "ok", Data: user, StatusCode: http.StatusOK}

	c.JSON(http.StatusOK, response)
}

func UpdateUserProfile(c *gin.Context) {
	requestContext := api.GetRequestContextFromRequest(c)
	span := api.AddTrace(requestContext, "info", "UpdateUserProfile")
	if span != nil {
		defer span.End()
	}

	log := requestContext.Log

	var user user_models.ResidentDTO
	err := c.ShouldBind(&user)
	if err != nil {
		log.Error("error while binding user profile update", zap.Error(err))
		c.Error(shared_errors.ErrInvalidPayload)
		return
	}

	err = user_facade.UpdateUserProfile(*requestContext, &user)
	if err != nil {
		log.Error("error while updating user profile", zap.Error(err))
		c.Error(err)
		return
	}
	response := shared_models.Response{Status: "ok", StatusCode: http.StatusOK}

	c.JSON(http.StatusOK, response)
}

func UpdateProfilePicture(c *gin.Context) {
	requestContext := api.GetRequestContextFromRequest(c)
	span := api.AddTrace(requestContext, "info", "UpdateProfilePicture")
	if span != nil {
		defer span.End()
	}

	log := requestContext.Log

	file, header, err := c.Request.FormFile("profile_picture")
	if err != nil {
		log.Error("error while getting file", zap.Error(err))
		c.Error(shared_errors.ErrInvalidPayload)
		return
	}

	err = user_facade.UpdateProfilePicture(*requestContext, file, header)
	if err != nil {
		log.Error("error while updating user profile picture", zap.Error(err))
		c.Error(err)
		return
	}

	response := shared_models.Response{Status: "ok", StatusCode: http.StatusOK}

	c.JSON(http.StatusOK, response)
}

func SendOTP(c *gin.Context) {
	requestContext := api.GetRequestContextFromRequest(c)
	span := api.AddTrace(requestContext, "info", "SendOTP")
	if span != nil {
		defer span.End()
	}

	log := requestContext.Log

	err := user_facade.SendOTP(*requestContext, c.Query("method"), c.Query("contact"))
	if err != nil {
		log.Error("error while sending otp", zap.Error(err))
		c.Error(err)
		return
	}
	response := shared_models.Response{Status: "ok", StatusCode: http.StatusOK}

	c.JSON(http.StatusOK, response)

}
