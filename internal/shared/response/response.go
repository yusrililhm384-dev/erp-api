package response

type Response struct {
	Success bool            `json:"success" example:"true"`
	Status  int             `json:"status" example:"200"`
	Message string          `json:"message" example:"Success"`
	Meta    *MetaPagination `json:"meta,omitempty"`
	Data    any             `json:"data,omitempty" example:"null"`
	Error   any             `json:"error,omitempty" example:"null"`
}

type MetaPagination struct {
	CurrentPage uint `json:"current_page"`
	Limit       uint `json:"limit"`
	TotalData   uint `json:"total_date"`
	TotalPages  uint `json:"total_pages"`
}
