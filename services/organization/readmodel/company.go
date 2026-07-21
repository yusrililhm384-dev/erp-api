package readmodel

import "time"

type CompanyListResponse struct {
	Meta *MetaPagination `json:"meta"`
	List []*CompanyList  `json:"list"`
}

type CompanyList struct {
	Id        uint      `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Website   *string   `json:"website,omitempty"`
	Location  *Point    `json:"location"`
	Country   string    `json:"country"`
	Language  string    `json:"language"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CompanyDetail struct {
	Id          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	Website     *string   `json:"website,omitempty"`
	Address     string    `json:"address"`
	Location    *Point    `json:"location"`
	Country     string    `json:"country"`
	Language    string    `json:"language"`
	Currency    string    `json:"currency"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
