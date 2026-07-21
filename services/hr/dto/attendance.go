package dto

import "time"

type SickReq struct {
	Note string `json:"note" validate:"required" example:"headeache"`
}

type CreateSickReq struct {
	Note          string `json:"note" validate:"required" example:"headeache"`
	AttachmentKey string `json:"attachment_key" validate:"required"`
}

type FilterDateReq struct {
	StartDate time.Time `json:"start_date" query:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" query:"end_date" validate:"required"`
}

type FilterDateStringReq struct {
	StartDate string `json:"start_date" query:"start_date" validate:"required"`
	EndDate   string `json:"end_date" query:"end_date" validate:"required"`
}

func (ds *FilterDateStringReq) DateStringToDate() (*FilterDateReq, error) {

	loc, err := time.LoadLocation("Asia/Jakarta")

	if err != nil {
		return nil, err
	}

	date := &FilterDateReq{}

	date.StartDate, err = time.ParseInLocation(time.DateOnly, ds.StartDate, loc)

	if err != nil {
		return nil, err
	}

	date.EndDate, err = time.ParseInLocation(time.DateOnly, ds.EndDate, loc)

	if err != nil {
		return nil, err
	}

	return date, nil
}
