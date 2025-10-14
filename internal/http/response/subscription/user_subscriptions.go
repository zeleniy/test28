package subscription_response

import "time"

type UserSubscription struct {
	ServiceName string    `boil:"service_name" json:"service_name"`
	Price       int       `boil:"price" json:"price"`
	UserUUID    string    `boil:"uuid" json:"user_id"`
	StartDate   time.Time `boil:"start_date" json:"start_date"`
}
