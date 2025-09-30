package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
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

	users, err := factory.CreateAndInsertUsers(ctx, db, 3,
		factory.UserLoginFunc(func() (string, error) { return faker.Username(), nil }),
		factory.UserPasswordHashFunc(func() (string, error) { return faker.Password(), nil }),
	)
	assert.NoError(t, err, "Failed to create users")

	subscriptionsCount := 0
	for _, user := range users {
		subscriptionsPerUser := rand.Intn(3) + 1
		_, err = factory.CreateAndInsertSubscriptions(ctx, db, subscriptionsPerUser,
			factory.SubscriptionWithUser(user),
			factory.SubscriptionServiceName([]string{"Okko", "Yandex", "Wink", "Sber", "Ivi"}[rand.Intn(5)]),
			factory.SubscriptionPrice([]int{10, 20, 30, 40, 50}[rand.Intn(5)]),
		)
		assert.NoError(t, err, "Failed to create subscriptions for user %d", user.ID)
		subscriptionsCount += subscriptionsPerUser
	}

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

	subscriptionsJson := json.Get("data.subscriptions")
	assert.True(t, subscriptionsJson.Exists())
	assert.True(t, subscriptionsJson.IsArray())
	assert.Len(t, subscriptionsJson.Array(), subscriptionsCount)

	subscriptionsJson.ForEach(func(_, subscription gjson.Result) bool {
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

func TestReadSubscription(t *testing.T) {

	user, err := factory.CreateAndInsertUser(ctx, db,
		factory.UserLogin(faker.Username()),
		factory.UserPasswordHash(faker.Password()),
	)

	assert.NoError(t, err, "Failed to create user")

	subscription, err := factory.CreateAndInsertSubscription(ctx, db,
		factory.SubscriptionUserID(user.ID),
		factory.SubscriptionServiceName("Ivi"),
		factory.SubscriptionPrice(100),
	)

	assert.NoError(t, err, "Failed to create subscription")

	req, err := http.NewRequest(http.MethodGet, "/subscriptions/"+strconv.Itoa(subscription.ID), nil)
	assert.NoError(t, err, "Failed to create request")

	w := httptest.NewRecorder()

	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code %d, got %d", http.StatusOK, w.Code)

	json := gjson.Parse(w.Body.String())

	if !json.IsObject() && !json.IsArray() {
		assert.True(t, json.Get("data").Exists(), "Response is not valid JSON")
	}

	assertResponseStructure(t, json)

	subscriptionJson := json.Get("data.subscription")
	assert.True(t, subscriptionJson.Exists())
	assert.True(t, subscriptionJson.IsObject())

	assert.Equal(t, int64(subscription.ID), subscriptionJson.Get("id").Int())
	assert.Equal(t, int64(subscription.UserID), subscriptionJson.Get("user_id").Int())
	assert.Equal(t, "Ivi", subscriptionJson.Get("service_name").Value())
	assert.Equal(t, int64(100), subscriptionJson.Get("price").Int())
	assert.IsType(t, "", subscriptionJson.Get("start_date").Value())
	assert.Nil(t, subscriptionJson.Get("end_date").Value())
	assert.IsType(t, "", subscriptionJson.Get("created_at").Value())
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
