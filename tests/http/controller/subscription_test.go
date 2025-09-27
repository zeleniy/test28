package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/zeleniy/test28/bootstrap"
	"github.com/zeleniy/test28/factory"
)

var (
	ginEngine *gin.Engine
	db        boil.ContextExecutor
	ctx       context.Context
)

func init() {

	ginEngine = bootstrap.SetUpApp(gin.TestMode, os.Getenv("DB_TEST_URL"))
	db = boil.GetContextDB()
	ctx = context.Background()
}

func TestGetSubscriptions(t *testing.T) {

	user, err := factory.CreateAndInsertUser(ctx, db,
		factory.UserLogin(faker.Username()),
		factory.UserPasswordHash(faker.Password()),
	)

	assert.NoError(t, err, "Failed to create user")

	_, err = factory.CreateSubscription(
		factory.SubscriptionUserID(user.ID),
		factory.SubscriptionServiceName("test-service"),
		factory.SubscriptionPrice(1000),
	)

	assert.NoError(t, err, "Failed to create subscription")

	req, err := http.NewRequest(http.MethodGet, "/subscriptions", nil)
	assert.NoError(t, err, "Failed to create request")

	w := httptest.NewRecorder()

	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code %d, got %d", http.StatusOK, w.Code)

	json := gjson.Parse(w.Body.String())

	if !json.IsObject() && !json.IsArray() {
		assert.True(t, json.Get("data").Exists(), "Response is not valid JSON")
	}

	assertResponseStructure(t, json)

	subscriptions := json.Get("data.subscriptions")
	assert.True(t, subscriptions.Exists())
	assert.True(t, subscriptions.IsArray())

	subscriptions.ForEach(func(_, subscription gjson.Result) bool {
		assert.Greater(t, subscription.Get("id").Int(), int64(0))
		assert.Greater(t, subscription.Get("user_id").Int(), int64(0))
		assert.NotEmpty(t, subscription.Get("service_name").String())
		assert.Greater(t, subscription.Get("price").Int(), int64(0))
		assert.IsType(t, "", subscription.Get("start_date").Value())
		assert.Nil(t, subscription.Get("end_date").Value())
		assert.IsType(t, "", subscription.Get("created_at").Value())
		return true
	})
}

func TestCreateSubscription(t *testing.T) {

	user, err := factory.CreateAndInsertUser(ctx, db,
		factory.UserLogin(faker.Username()),
		factory.UserPasswordHash(faker.Password()),
	)

	assert.NoError(t, err, "Failed to create user")

	data := map[string]interface{}{
		"user_id":      user.ID,
		"service_name": "Okko",
		"price":        100,
	}

	jsonData, err := json.Marshal(data)

	if err != nil {
		assert.NoError(t, err, "Failed to create JSON")
	}

	req, err := http.NewRequest(http.MethodPost, "/subscriptions", bytes.NewBuffer(jsonData))
	assert.NoError(t, err, "Failed to create request")

	w := httptest.NewRecorder()

	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code %d, got %d", http.StatusOK, w.Code)

	json := gjson.Parse(w.Body.String())

	if !json.IsObject() && !json.IsArray() {
		assert.True(t, json.Get("data").Exists(), "Response is not valid JSON")
	}

	assertResponseStructure(t, json)

	subscription := json.Get("data.subscription")
	assert.True(t, subscription.Exists(), "Response does not contain 'data.subscription' key")

	assert.Greater(t, subscription.Get("id").Int(), int64(0))
	assert.Equal(t, subscription.Get("user_id").Int(), int64(user.ID))
	assert.Equal(t, subscription.Get("service_name").String(), "Okko")
	assert.Equal(t, subscription.Get("price").Int(), int64(100))
	assert.IsType(t, "", subscription.Get("start_date").Value())
	assert.Nil(t, subscription.Get("end_date").Value())
	assert.IsType(t, "", subscription.Get("created_at").Value())
}

func assertResponseStructure(t *testing.T, json gjson.Result) {

	assert.True(t, json.Get("data").Exists(), "Response does not contain 'data' key")
	assert.True(t, json.Get("meta").Exists(), "Response does not contain 'meta' key")
	assert.True(t, json.Get("error").Exists(), "Response does not contain 'error' key")

	timestamp := json.Get("meta.timestamp")
	assert.True(t, timestamp.Exists(), "Response does not contain 'timestamp' in 'meta'")
	assert.NotEmpty(t, timestamp.String(), "Timestamp should not be empty")

	assert.Nil(t, json.Get("error").Value(), "Error should be empty")
}
