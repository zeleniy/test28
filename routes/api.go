package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeleniy/test28/http/controllers"
)

func SetupRoutes(ginEngine *gin.Engine) {

	subscriptionCtrl := &controllers.SubscriptionController{}

	ginEngine.GET("/ping", func(ginContext *gin.Context) {
		ginContext.Header("Content-Type", "text/plain")
		ginContext.String(http.StatusOK, "pong")
	})

	subscriptions := ginEngine.Group("/subscriptions")

	subscriptions.GET("", subscriptionCtrl.GetSubscriptions)
	subscriptions.POST("", subscriptionCtrl.CreateSubscription)
	subscriptions.GET("/:id", subscriptionCtrl.ReadSubscription)
	subscriptions.PATCH("/:id", subscriptionCtrl.UpdateSubscription)
	subscriptions.PUT("/:id", subscriptionCtrl.UpdateSubscription)
	subscriptions.DELETE("/:id", subscriptionCtrl.DeleteSubscription)
	subscriptions.POST("/report", subscriptionCtrl.GetAccountingReport)
}
