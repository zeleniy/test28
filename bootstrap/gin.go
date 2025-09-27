package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/zeleniy/test28/http/middleware"
	"github.com/zeleniy/test28/routes"
)

func SetUpGin(ginMode string) *gin.Engine {

	gin.SetMode(ginMode)

	gin := gin.Default()

	gin.Use(middleware.DataWrapperMiddleware())

	routes.SetupRoutes(gin)

	return gin
}
