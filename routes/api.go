package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zeleniy/test28/http/controllers"
)

func SetupRoutes(r *gin.Engine) {

	subscriptionCtrl := &controllers.SubscriptionController{}

	r.GET("/subscriptions", subscriptionCtrl.GetSubscriptions)
	r.POST("/subscriptions", subscriptionCtrl.CreateSubscription)
}
