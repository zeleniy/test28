package controllers

import (

	// "time"

	"context"
	"net/http"

	"github.com/aarondl/sqlboiler/v4/boil"
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

	var input struct {
		UserID      int    `json:"user_id" binding:"required,gt=1"`
		ServiceName string `json:"service_name" binding:"required"`
		Price       int    `json:"price" binding:"required,gt=1"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription := models.Subscription{
		UserID:      input.UserID,
		ServiceName: input.ServiceName,
		Price:       input.Price,
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
