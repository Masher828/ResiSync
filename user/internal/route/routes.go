package routes

import (
	"ResiSync/user/internal/api/controller"

	"github.com/gin-gonic/gin"
)

type Rest struct{}

func (r *Rest) SetupPrivateRoutes(engine *gin.Engine) {

}

func (r *Rest) SetupPublicRoutes(engine *gin.Engine) {
	engine.GET("/user", controller.GetUserDetails)
}
