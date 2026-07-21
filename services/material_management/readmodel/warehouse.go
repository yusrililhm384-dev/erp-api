package readmodel

import "time"

type WarehouseListResponse struct {
	Meta *MetaPagination  `json:"meta"`
	List []*WarehouseList `json:"list"`
}

type WarehouseList struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	Location    *Point    `json:"location"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Point struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
