package readmodel

import "time"

type PositionCategoryListResponse struct {
	Meta *MetaPagination         `json:"meta"`
	List []*PositionCategoryList `json:"list"`
}

type PositionCategoryDetail struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PositionCategoryList struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
