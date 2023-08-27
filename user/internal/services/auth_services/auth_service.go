package authservices

import (
	"ResiSync/pkg/api"
	"ResiSync/pkg/models"
	authrepository "ResiSync/user/internal/repository/auth_repository"
	"fmt"
)

func GetUserDetailsService(requestContext *models.ResiSyncRequestContext, id int64) {
	span := api.AddTrace(requestContext, "info", "GetUserDetailsService")
	if span != nil {
		defer span.End()
	}
	fmt.Println("Hello there")

	authrepository.GetUserDetailsFromRepo(requestContext)
}
