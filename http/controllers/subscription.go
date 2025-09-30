package controllers

import (

	// "time"

	"context"
	"net/http"

	"github.com/aarondl/sqlboiler/v4/boil"
	requests "github.com/zeleniy/test28/http/requests/subscription"
	"github.com/zeleniy/test28/models"

	"github.com/gin-gonic/gin"
)

type SubscriptionController struct{}

// Get subscriptions
func (ctrl *SubscriptionController) GetSubscriptions(c *gin.Context) {

	subscriptions, err := models.Subscriptions().All(context.Background(), boil.GetContextDB())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Set("data", map[string]interface{}{
		"subscriptions": subscriptions,
	})
}

// Subscribe user
func (ctrl *SubscriptionController) CreateSubscription(c *gin.Context) {

	var request requests.CreateSubscriptionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription := models.Subscription{
		UserID:      request.UserID,
		ServiceName: request.ServiceName,
		Price:       request.Price,
	}

	err := subscription.Insert(c.Request.Context(), boil.GetContextDB(), boil.Infer())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Set("data", map[string]interface{}{
		"subscription": subscription,
	})
}

// Get user's subscription info
func (ctrl *SubscriptionController) ReadSubscription(c *gin.Context) {

	var request requests.ReadSubscriptionRequest

	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := models.FindSubscription(context.Background(), boil.GetContextDB(), request.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Set("data", map[string]interface{}{
		"subscription": subscription,
	})
}
