package request

type UUIDRequest struct {
	UUID int `uri:"id" binding:"required,len=36"`
}
