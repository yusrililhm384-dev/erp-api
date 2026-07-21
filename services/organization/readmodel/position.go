package readmodel

import "time"

type PositionListResponse struct {
	Meta *MetaPagination `json:"meta"`
	List []*PositionList `json:"list"`
}

type PositionList struct {
	Id         uint      `json:"id"`
	Name       string    `json:"name"`
	ApprovedBy string    `json:"approved_by"`
	Category   string    `json:"category"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PositionDetail struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ApprovedBy  string    `json:"approved_by"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
