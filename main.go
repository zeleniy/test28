package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/zeleniy/test28/bootstrap"
)

func main() {

	_, err := bootstrap.SetUpDb(os.Getenv("DB_URL"))

	if err != nil {
		panic(err)
	}

	r := bootstrap.SetUpGin(gin.ReleaseMode)

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
