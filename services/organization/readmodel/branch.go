package readmodel

import "time"

type BranchListResponse struct {
	Meta *MetaPagination `json:"meta"`
	List []*BranchList   `json:"list"`
}

type BranchDetail struct {
	Id          uint      `json:"id"`
	Company     string    `json:"company"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	Location    *Point    `json:"location"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BranchList struct {
	Id          uint      `json:"id"`
	Company     string    `json:"company"`
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
