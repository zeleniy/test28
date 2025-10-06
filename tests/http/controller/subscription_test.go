package controller

import (
	"bytes"
	"context"
	"database/sql"
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

func withTransaction(t *testing.T, testFunc func(tx *sql.Tx)) {

	tx, err := db.(*sql.DB).Begin()
	if err != nil {
		t.Fatalf("Cannot begin transaction: %v", err)
	}
	defer tx.Rollback()

	originalDB := boil.GetDB()
	boil.SetDB(tx)
	defer boil.SetDB(originalDB)

	testFunc(tx)
}

func TestGetSubscriptions(t *testing.T) {

	withTransaction(t, func(tx *sql.Tx) {

		users, err := factory.CreateAndInsertUsers(ctx, tx, 3,
			factory.UserLoginFunc(func() (string, error) { return faker.Username(), nil }),
			factory.UserPasswordHashFunc(func() (string, error) { return faker.Password(), nil }),
		)
		assert.NoError(t, err, "Failed to create users")

		subscriptionsCount := 0
		for _, user := range users {
			subscriptionsPerUser := rand.Intn(3) + 1
			_, err = factory.CreateAndInsertSubscriptions(ctx, tx, subscriptionsPerUser,
				factory.SubscriptionWithUser(user),
				factory.SubscriptionServiceName([]string{"Okko", "Yandex", "Wink", "Sber", "Ivi"}[rand.Intn(5)]),
				factory.SubscriptionPrice([]int{10, 20, 30, 40, 50}[rand.Intn(5)]),
			)
			assert.NoError(t, err, "Failed to create subscriptions for user %d", user.ID)
			subscriptionsCount += subscriptionsPerUser
		}

		gjsonBody := sendAndTestRequest(t, http.MethodGet, "/subscriptions", http.StatusOK, nil)
		assertResponseStructure(t, gjsonBody)

		gjsonSubscriptions := gjsonBody.Get("data.subscriptions")
		assert.True(t, gjsonSubscriptions.Exists())
		assert.True(t, gjsonSubscriptions.IsArray())
		assert.Len(t, gjsonSubscriptions.Array(), subscriptionsCount)

		gjsonSubscriptions.ForEach(func(_, subscription gjson.Result) bool {
			assert.Greater(t, subscription.Get("id").Int(), int64(0))
			assert.Greater(t, subscription.Get("user_id").Int(), int64(0))
			assert.NotEmpty(t, subscription.Get("service_name").String())
			assert.Greater(t, subscription.Get("price").Int(), int64(0))
			assert.IsType(t, "", subscription.Get("start_date").Value())
			assert.Nil(t, subscription.Get("end_date").Value())
			assert.IsType(t, "", subscription.Get("created_at").Value())
			return true
		})
	})
}

func TestGetAccountingReport(t *testing.T) {

	withTransaction(t, func(tx *sql.Tx) {

		user, err := factory.CreateAndInsertUser(ctx, tx,
			factory.UserLoginFunc(func() (string, error) { return faker.Username(), nil }),
			factory.UserPasswordHashFunc(func() (string, error) { return faker.Password(), nil }),
		)
		assert.NoError(t, err, "Failed to create users")

		_, err = factory.CreateAndInsertSubscription(ctx, tx,
			factory.SubscriptionWithUser(user),
			factory.SubscriptionServiceName("Yandex"),
			factory.SubscriptionPrice(10),
		)
		assert.NoError(t, err, "Failed to create subscriptions for user %d", user.ID)

		_, err = factory.CreateAndInsertSubscription(ctx, tx,
			factory.SubscriptionWithUser(user),
			factory.SubscriptionServiceName("Okko"),
			factory.SubscriptionPrice(20),
		)
		assert.NoError(t, err, "Failed to create subscriptions for user %d", user.ID)

		gjsonBody := sendAndTestRequest(t, http.MethodPost, "/subscriptions/report", http.StatusOK, map[string]interface{}{})
		assertResponseStructure(t, gjsonBody)

		assert.Equal(t, gjsonBody.Get("data.count").Int(), int64(2))
		assert.Equal(t, gjsonBody.Get("data.sum").Int(), int64(30))

		// Test with another one user

		user, err = factory.CreateAndInsertUser(ctx, tx,
			factory.UserLoginFunc(func() (string, error) { return faker.Username(), nil }),
			factory.UserPasswordHashFunc(func() (string, error) { return faker.Password(), nil }),
		)
		assert.NoError(t, err, "Failed to create users")

		_, err = factory.CreateAndInsertSubscription(ctx, tx,
			factory.SubscriptionWithUser(user),
			factory.SubscriptionServiceName("Ivi"),
			factory.SubscriptionPrice(10),
		)
		assert.NoError(t, err, "Failed to create subscriptions for user %d", user.ID)

		gjsonBody = sendAndTestRequest(t, http.MethodPost, "/subscriptions/report", http.StatusOK, map[string]interface{}{
			"user_id": user.ID,
		})

		assertResponseStructure(t, gjsonBody)
		assert.Equal(t, gjsonBody.Get("data.count").Int(), int64(1))
		assert.Equal(t, gjsonBody.Get("data.sum").Int(), int64(10))

		gjsonBody = sendAndTestRequest(t, http.MethodPost, "/subscriptions/report", http.StatusOK, map[string]interface{}{
			"user_id":      user.ID,
			"service_name": "Ivi",
			"to":           "01-01-2025",
		})

		assertResponseStructure(t, gjsonBody)
		assert.Equal(t, gjsonBody.Get("data.count").Int(), int64(1))
		assert.Equal(t, gjsonBody.Get("data.sum").Int(), int64(10))

		gjsonBody = sendAndTestRequest(t, http.MethodPost, "/subscriptions/report", http.StatusOK, map[string]interface{}{
			"user_id":      user.ID,
			"service_name": "Okko",
			"from":         "01-01-1901",
		})

		assertResponseStructure(t, gjsonBody)
		assert.Equal(t, gjsonBody.Get("data.count").Int(), int64(0))
		assert.Equal(t, gjsonBody.Get("data.sum").Int(), int64(0))
	})
}

func TestCreateSubscription(t *testing.T) {

	withTransaction(t, func(tx *sql.Tx) {

		user, err := factory.CreateAndInsertUser(ctx, tx,
			factory.UserLogin(faker.Username()),
			factory.UserPasswordHash(faker.Password()),
		)

		assert.NoError(t, err, "Failed to create user")

		gjsonBody := sendAndTestRequest(t, http.MethodPost, "/subscriptions", http.StatusOK, map[string]interface{}{
			"user_id":      user.ID,
			"service_name": "Okko",
			"price":        100,
		})

		assertResponseStructure(t, gjsonBody)
		gjsonSubscription := gjsonBody.Get("data.subscription")
		assert.True(t, gjsonSubscription.Exists(), "Response does not contain 'data.subscription' key")

		assert.Greater(t, gjsonSubscription.Get("id").Int(), int64(0))
		assert.Equal(t, gjsonSubscription.Get("user_id").Int(), int64(user.ID))
		assert.Equal(t, gjsonSubscription.Get("service_name").String(), "Okko")
		assert.Equal(t, gjsonSubscription.Get("price").Int(), int64(100))
		assert.IsType(t, "", gjsonSubscription.Get("start_date").Value())
		assert.Nil(t, gjsonSubscription.Get("end_date").Value())
		assert.IsType(t, "", gjsonSubscription.Get("created_at").Value())
	})
}

func TestReadSubscription(t *testing.T) {

	withTransaction(t, func(tx *sql.Tx) {

		user, err := factory.CreateAndInsertUser(ctx, tx,
			factory.UserLogin(faker.Username()),
			factory.UserPasswordHash(faker.Password()),
		)

		assert.NoError(t, err, "Failed to create user")

		subscription, err := factory.CreateAndInsertSubscription(ctx, tx,
			factory.SubscriptionUserID(user.ID),
			factory.SubscriptionServiceName("Ivi"),
			factory.SubscriptionPrice(100),
		)

		assert.NoError(t, err, "Failed to create subscription")

		gjsonBody := sendAndTestRequest(t, http.MethodGet, "/subscriptions/"+strconv.Itoa(subscription.ID), http.StatusOK, nil)
		assertResponseStructure(t, gjsonBody)

		gjsonSubscription := gjsonBody.Get("data.subscription")
		assert.True(t, gjsonSubscription.Exists())
		assert.True(t, gjsonSubscription.IsObject())
		assert.Equal(t, int64(subscription.ID), gjsonSubscription.Get("id").Int())
		assert.Equal(t, int64(subscription.UserID), gjsonSubscription.Get("user_id").Int())
		assert.Equal(t, "Ivi", gjsonSubscription.Get("service_name").Value())
		assert.Equal(t, int64(100), gjsonSubscription.Get("price").Int())
		assert.IsType(t, "", gjsonSubscription.Get("start_date").Value())
		assert.Nil(t, gjsonSubscription.Get("end_date").Value())
		assert.IsType(t, "", gjsonSubscription.Get("created_at").Value())
	})
}

func TestUpdateSubscription(t *testing.T) {

	withTransaction(t, func(tx *sql.Tx) {

		user, err := factory.CreateAndInsertUser(ctx, tx,
			factory.UserLogin(faker.Username()),
			factory.UserPasswordHash(faker.Password()),
		)

		assert.NoError(t, err, "Failed to create user")

		subscription, err := factory.CreateAndInsertSubscription(ctx, tx,
			factory.SubscriptionUserID(user.ID),
			factory.SubscriptionServiceName("Ivi"),
			factory.SubscriptionPrice(100),
		)

		assert.NoError(t, err, "Failed to create subscription")

		httpMethod := []string{http.MethodPatch, http.MethodPut}[rand.Intn(2)]
		sendAndTestRequest(t, httpMethod, "/subscriptions/"+strconv.Itoa(subscription.ID), http.StatusMethodNotAllowed, nil)
	})
}

func TestDeleteSubscription(t *testing.T) {

	withTransaction(t, func(tx *sql.Tx) {

		user, err := factory.CreateAndInsertUser(ctx, tx,
			factory.UserLogin(faker.Username()),
			factory.UserPasswordHash(faker.Password()),
		)

		assert.NoError(t, err, "Failed to create user")

		subscription, err := factory.CreateAndInsertSubscription(ctx, tx,
			factory.SubscriptionUserID(user.ID),
			factory.SubscriptionServiceName("Ivi"),
			factory.SubscriptionPrice(100),
		)

		assert.NoError(t, err, "Failed to create subscription")

		sendAndTestRequest(t, http.MethodDelete, "/subscriptions/"+strconv.Itoa(subscription.ID), http.StatusNoContent, nil)
	})
}

func sendAndTestRequest(t *testing.T, httpMethod, url string, code int, data map[string]interface{}) gjson.Result {

	jsonData, err := json.Marshal(data)

	if err != nil {
		assert.NoError(t, err, "Failed to create JSON")
	}

	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonData))
	assert.NoError(t, err, "Failed to create request")
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)
	assert.Equal(t, code, w.Code, "Expected status code %d, got %d", code, w.Code)
	gjsonBody := gjson.Parse(w.Body.String())

	return gjsonBody
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
