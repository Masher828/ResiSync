package shared_api

import (
	shared_errors "ResiSync/shared/errors"
	shared_models "ResiSync/shared/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			statusCode := http.StatusInternalServerError
			resp := shared_models.Response{
				Status:     "error",
				Error:      "internal server error",
				StatusCode: statusCode,
			}
			if shared_errors.IsInternalError(c.Errors[0]) {
				resp.Error = c.Errors[0].Error()
				resp.StatusCode = http.StatusBadRequest
			}
			c.JSON(statusCode, resp)
		}
	}
}
