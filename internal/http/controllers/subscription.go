package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/gin-gonic/gin"
	"github.com/zeleniy/test28/internal/http/request"
	subscription_request "github.com/zeleniy/test28/internal/http/request/subscription"
	subscription_response "github.com/zeleniy/test28/internal/http/response/subscription"
	"github.com/zeleniy/test28/internal/models"
)

type SubscriptionController struct{}

// Get subscriptions
func (ctrl *SubscriptionController) GetSubscriptions(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var subscriptions []subscription_response.UserSubscription

	err := models.NewQuery(
		qm.Select("service_name", "price", "uuid", "start_date"),
		qm.From("subscriptions"),
		qm.InnerJoin("users on subscriptions.user_id = users.id"),
	).Bind(ctx, boil.GetContextDB(), &subscriptions)

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

	var request subscription_request.CreateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.Users(qm.Where("uuid=?", request.UserUUID)).One(c.Request.Context(), boil.GetContextDB())

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	subscription := models.Subscription{
		UserID:      user.ID,
		ServiceName: request.ServiceName,
		Price:       request.Price,
	}

	err = subscription.Insert(c.Request.Context(), boil.GetContextDB(), boil.Infer())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userSubscription := subscription_response.UserSubscription{
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		UserUUID:    user.UUID,
		StartDate:   subscription.StartDate,
	}

	c.Set("data", map[string]interface{}{
		"subscription": userSubscription,
	})
}

// Get user's subscription info
func (ctrl *SubscriptionController) ReadSubscription(c *gin.Context) {

	var request request.IdRequest

	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var subscription subscription_response.UserSubscription

	err := models.NewQuery(
		qm.Select("service_name", "price", "uuid", "start_date"),
		qm.From("subscriptions"),
		qm.InnerJoin("users on subscriptions.user_id = users.id"),
		qm.Where("subscriptions.id = ?", request.ID),
	).Bind(context.Background(), boil.GetContextDB(), &subscription)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Set("data", map[string]interface{}{
		"subscription": subscription,
	})
}

// Update user's subscription details
func (ctrl *SubscriptionController) UpdateSubscription(c *gin.Context) {

	c.AbortWithStatus(http.StatusMethodNotAllowed)
}

// Cancel subscription
func (ctrl *SubscriptionController) DeleteSubscription(c *gin.Context) {

	var request request.IdRequest

	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := models.Subscriptions(models.SubscriptionWhere.ID.EQ(request.ID)).
		DeleteAll(context.Background(), boil.GetContextDB())

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Get accounting report
func (ctrl *SubscriptionController) GetAccountingReport(c *gin.Context) {
	boil.DebugMode = true
	boil.DebugWriter = os.Stdout // или твой логгер

	var request subscription_request.ReportRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mods, err := getGetAccountingReportCriteria(request, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mods = append(mods, qm.Select("COUNT(*) as count, COALESCE(SUM(price), 0) as sum"))

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var result struct {
		Count int `boil:"count"`
		Sum   int `boil:"sum"`
	}

	err = models.Subscriptions(mods...).Bind(ctx, boil.GetContextDB(), &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Set("data", map[string]interface{}{
		"sum":   result.Sum,
		"count": result.Count,
		"from":  request.From,
		"to":    request.To,
	})
}

func getGetAccountingReportCriteria(request subscription_request.ReportRequest, c *gin.Context) ([]qm.QueryMod, error) {

	var mods []qm.QueryMod

	if request.From != nil {
		if fromDate, err := time.Parse("02-01-2006", *request.From); err != nil {
			return nil, err
		} else {
			mods = append(mods, models.SubscriptionWhere.StartDate.GTE(fromDate))
		}
	}

	if request.To != nil {
		if toDate, err := time.Parse("02-01-2006", *request.To); err != nil {
			return nil, err
		} else {
			mods = append(mods, models.SubscriptionWhere.EndDate.LTE(null.TimeFrom(toDate)))
		}
	}

	if request.UserUUID != nil {
		mods = append(mods, models.UserWhere.UUID.EQ(*request.UserUUID), qm.InnerJoin("users on subscriptions.user_id = users.id"))
	}

	if request.ServiceName != nil {
		mods = append(mods, models.SubscriptionWhere.ServiceName.EQ(*request.ServiceName))
	}

	return mods, nil
}
