package readmodel

import "time"

type MaterialTypeListResponse struct {
	Meta *MetaPagination     `json:"meta"`
	List []*MaterialTypeList `json:"list"`
}

type MaterialTypeList struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MaterialTypeDetail struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
