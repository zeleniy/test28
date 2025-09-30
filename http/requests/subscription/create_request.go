package requests

type CreateSubscriptionRequest struct {
	UserID      int    `json:"user_id" binding:"required,gt=1"`
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"required,gt=1"`
}
