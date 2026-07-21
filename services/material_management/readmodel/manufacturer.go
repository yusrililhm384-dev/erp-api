package readmodel

import "time"

type ManufacturerListResponse struct {
	Meta *MetaPagination     `json:"meta"`
	List []*ManufacturerList `json:"list"`
}

type ManufacturerList struct {
	Id        uint      `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ManufacturerDetail struct {
	Id        uint      `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
