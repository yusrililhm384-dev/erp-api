package readmodel

import "time"

type PaymentTermListResponse struct {
	Meta *MetaPagination    `json:"meta"`
	List []*PaymentTermList `json:"list"`
}

type PaymentTermList struct {
	Id          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PaymentTermDetail struct {
	Id          uint             `json:"id"`
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Countries   []*CountryDetail `json:"coutries"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type PaymentMethodListResponse struct {
	Meta *MetaPagination      `json:"meta"`
	List []*PaymentMethodList `json:"list"`
}

type PaymentMethodList struct {
	Id          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PaymentMethodDetail struct {
	Id          uint             `json:"id"`
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Countries   []*CountryDetail `json:"coutries"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}
