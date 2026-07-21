package readmodel

import "time"

type CountryListResponse struct {
	Meta *MetaPagination `json:"meta"`
	List []*CountryList  `json:"list"`
}

type CountryList struct {
	Id        uint      `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CountryDetail struct {
	Id        uint      `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
}
