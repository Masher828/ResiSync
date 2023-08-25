package shared_api

import "github.com/gin-gonic/gin"

func HandleError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
