package readmodel

import "time"

type CurrencyListResponse struct {
	Meta *MetaPagination `json:"meta"`
	List []*CurrencyList `json:"list"`
}

type CurrencyList struct {
	Id        uint      `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CurrencyDetail struct {
	Id        uint             `json:"id"`
	Code      string           `json:"code"`
	Name      string           `json:"name"`
	Countries []*CountryDetail `json:"coutries"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}
