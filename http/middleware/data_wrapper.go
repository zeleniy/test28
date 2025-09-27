package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func DataWrapperMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()

		data, exists := c.Get("data")
		if !exists {
			return
		}

		response := map[string]interface{}{
			"data": data,
			"meta": map[string]interface{}{
				"timestamp": time.Now(),
			},
			"error": nil,
		}

		c.JSON(http.StatusOK, response)
	}
}
