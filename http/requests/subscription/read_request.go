package requests

type ReadSubscriptionRequest struct {
	ID int `uri:"id" binding:"required,gt=1"`
}
