package controller

import (
	"ResiSync/pkg/api"
	authservices "ResiSync/user/internal/services/auth_services"

	"github.com/gin-gonic/gin"
)

func GetUserDetails(c *gin.Context) {
	requestContext := api.GetRequestContextFromRequest(c)
	span := api.AddTrace(requestContext, "info", "GetUserDetails")
	if span != nil {
		defer span.End()
	}

	authservices.GetUserDetailsService(requestContext, 123)
}
