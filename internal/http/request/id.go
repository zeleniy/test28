package request

type IdRequest struct {
	ID int `uri:"id" binding:"required,gt=0"`
}
