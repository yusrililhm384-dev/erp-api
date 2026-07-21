package readmodel

import "time"

type MaterialGroupListResponse struct {
	Meta *MetaPagination      `json:"meta"`
	List []*MaterialGroupList `json:"list"`
}

type MaterialGroupList struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MaterialGroupDetail struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
