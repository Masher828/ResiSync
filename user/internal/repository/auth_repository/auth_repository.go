package authrepository

import (
	"ResiSync/pkg/api"
	"ResiSync/pkg/models"
	userModels "ResiSync/user/internal/models"
	"fmt"
)

func GetUserDetailsFromRepo(requestContext *models.ResiSyncRequestContext) {
	span := api.AddTrace(requestContext, "info", "GetUserDetailsFromRepo")
	if span != nil {
		defer span.End()
	}
	fmt.Println("Getting data from repo")

	user := userModels.User{}

	fmt.Println(user)

}
