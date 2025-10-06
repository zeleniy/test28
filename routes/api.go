package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zeleniy/test28/http/controllers"
)

func SetupRoutes(r *gin.Engine) {

	subscriptionCtrl := &controllers.SubscriptionController{}

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/subscriptions", subscriptionCtrl.GetSubscriptions)
	r.POST("/subscriptions", subscriptionCtrl.CreateSubscription)
	r.GET("/subscriptions/:id", subscriptionCtrl.ReadSubscription)
	r.PATCH("/subscriptions/:id", subscriptionCtrl.UpdateSubscription)
	r.PUT("/subscriptions/:id", subscriptionCtrl.UpdateSubscription)
	r.DELETE("/subscriptions/:id", subscriptionCtrl.DeleteSubscription)
	r.POST("/subscriptions/report", subscriptionCtrl.GetAccountingReport)
}
