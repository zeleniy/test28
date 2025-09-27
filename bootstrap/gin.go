package bootstrap

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeleniy/test28/http/middleware"
	"github.com/zeleniy/test28/routes"
)

func SetUpGin(ginMode string) *gin.Engine {

	gin.SetMode(ginMode)

	r := gin.Default()

	r.Use(middleware.DataWrapperMiddleware())

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	// r.GET("/user/:name", func(c *gin.Context) {
	// 	user := c.Params.ByName("name")
	// 	value, ok := db[user]
	// 	if ok {
	// 		c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
	// 	} else {
	// 		c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
	// 	}
	// })

	routes.SetupRoutes(r)

	return r
}
