package readmodel

type MetaPagination struct {
	CurrentPage uint `json:"current_page"`
	Limit       uint `json:"limit"`
	TotalData   uint `json:"total_date"`
	TotalPages  uint `json:"total_pages"`
}
