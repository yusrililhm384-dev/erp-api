package dto

type CreateCompanyReq struct {
	Code        string  `json:"code" example:"xiaoching" validate:"required"`
	Name        string  `json:"name" example:"Xiaoching Teknologi" validate:"required"`
	Description string  `json:"description" example:"PT Xiaoching Teknologi" validate:"required"`
	Phone       string  `json:"phone" example:"xxxxxxxxxxxx" validate:"required"`
	Email       string  `json:"email" example:"test@xiaoching.com" validate:"required"`
	Website     *string `json:"website,omitempty" example:"xiaoching.com" validate:"required"`
	Address     string  `json:"address" example:"South Jakarta" validate:"required"`
	Location    *Point  `json:"location" validate:"required"`
	CountryId   uint    `json:"country_id" example:"1" validate:"required"`
	CurrencyId  uint    `json:"currency_id" example:"1" validate:"required"`
	LanguageId  uint    `json:"language_id" example:"1" validate:"required"`
}

type UpdateCompanyReq struct {
	Code        string  `json:"code" example:"xiaoching"`
	Name        string  `json:"name" example:"Xiaoching Teknologi"`
	Description string  `json:"description" example:"PT Xiaoching Teknologi"`
	Phone       string  `json:"phone" example:"xxxxxxxxxxxx"`
	Email       string  `json:"email" example:"test@xiaoching.com"`
	Website     *string `json:"website,omitempty" example:"xiaoching.com"`
	Address     string  `json:"address" example:"South Jakarta"`
	Location    *Point  `json:"location"`
	CountryId   uint    `json:"country_id" example:"1"`
	CurrencyId  uint    `json:"currency_id" example:"1"`
	LanguageId  uint    `json:"language_id" example:"1"`
}
