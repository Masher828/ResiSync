package routes

import (
	"ResiSync/app/internal/api/controller"

	"github.com/gin-gonic/gin"
)

type Rest struct{}

func (r *Rest) SetupPrivateRoutes(engine *gin.Engine) {

}

func (r *Rest) SetupPublicRoutes(engine *gin.Engine) {

	//user Auth
	engine.POST("/user/signin", controller.SignIn)
	engine.POST("/user/signup", controller.SignUp)
	engine.POST("/user/logout", controller.LogOut)

	//user Details
	engine.GET("/user/profile", controller.GetUserProfile)
	engine.PUT("/user/profile", controller.UpdateUserProfile)
	engine.PUT("/user/profile/picture", controller.UpdateProfilePicture)
}
