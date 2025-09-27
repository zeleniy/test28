package bootstrap

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func SetUpApp(ginMode string, dsn string) *gin.Engine {

	_, err := SetUpDb(dsn)

	if err != nil {
		panic(err)
	}

	return SetUpGin(gin.ReleaseMode)
}
