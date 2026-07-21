package readmodel

import "time"

type IndustrySectorListResponse struct {
	Meta *MetaPagination       `json:"meta"`
	List []*IndustrySectorList `json:"list"`
}

type IndustrySectorList struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type IndustrySectorDetail struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
