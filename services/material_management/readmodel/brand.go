package readmodel

import "time"

type BrandListResponse struct {
	Meta *MetaPagination `json:"meta"`
	List []*BrandList    `json:"list"`
}

type BrandList struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BrandDetail struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
