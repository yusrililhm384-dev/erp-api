package readmodel

import "time"

type MaterialMasterListResponse struct {
	Meta *MetaPagination       `json:"meta"`
	List []*MaterialMasterList `json:"list"`
}

type MaterialMasterList struct {
	Id            uint      `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	MaterialType  string    `json:"material_type"`
	MaterialGroup string    `json:"material_group"`
	BaseUnit      string    `json:"base_unit"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type MaterialMasterDetail struct {
	Id        uint      `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
