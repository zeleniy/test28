package subscription_request

type CreateRequest struct {
	UserUUID    string `json:"user_id" binding:"required,len=36"`
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"required,gt=0"`
}
