package subscription

type ReportRequest struct {
	UserID      *int    `json:"user_id" binding:"omitempty,gt=0"`
	ServiceName *string `json:"service_name" binding:"omitempty,min=1,max=255"`
	From        *string `json:"from_date" binding:"omitempty,regex=^\\d{2}-\\d{2}-\\d{4}$,date=02-01-2006"`
	To          *string `json:"to_date" binding:"omitempty,regex=^\\d{2}-\\d{2}-\\d{4}$,date=02-01-2006"`
}
