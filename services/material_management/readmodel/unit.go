package readmodel

import "time"

type MaterialUnitListResponse struct {
	Meta *MetaPagination     `json:"meta"`
	List []*MaterialUnitList `json:"list"`
}

type MaterialUnitList struct {
	Id        uint      `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MaterialUnitDetail struct {
	Id          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MaterialUnitCategoryListResponse struct {
	Meta *MetaPagination             `json:"meta"`
	List []*MaterialUnitCategoryList `json:"list"`
}

type MaterialUnitCategoryList struct {
	Id        uint      `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MaterialUnitCategoryDetail struct {
	Id        uint            `json:"id"`
	Code      string          `json:"code"`
	Name      string          `json:"name"`
	List      []*MaterialUnit `json:"list"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type MaterialUnit struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
