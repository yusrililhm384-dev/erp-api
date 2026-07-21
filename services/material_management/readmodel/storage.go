package readmodel

import "time"

type StorageLocationListResponse struct {
	Meta *MetaPagination        `json:"meta"`
	List []*StorageLocationList `json:"list"`
}

type StorageLocationList struct {
	Id        uint      `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StorageLocationDetail struct {
	Id          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
