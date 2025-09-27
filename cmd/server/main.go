package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/zeleniy/test28/bootstrap"
)

func main() {

	bootstrap.SetUpApp(gin.ReleaseMode, os.Getenv("DB_URL")).
		Run(":8080")
}
