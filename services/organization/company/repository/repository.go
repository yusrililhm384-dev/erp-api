package repository

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"enterprise_resource_planning/services/organization/entity"
	"enterprise_resource_planning/services/organization/readmodel"

	errs "enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

const (
	companyListQuery = `
		select
			c.id,
			c.code,
			c.name,
			c.phone,
			c.email,
			c.website,
			ST_X(c.location::geometry) AS longitude,
			ST_Y(c.location::geometry) AS latitude,
			cty.name as country,
			lang.name as language,
			cry.name as currency,
			c.created_at,
			c.updated_at
		from companies c
		inner join countries cty on c.country_id = cty.id and cty.deleted_at is null 
		inner join languages lang on c.language_id = lang.id and lang.deleted_at is null
		inner join currencies cry on c.currency_id = cry.id and cry.deleted_at is null
		where c.deleted_at is null
		order by c.created_at desc
		limit $1 offset $2
	`

	countCompanyListQuery = `
		select count(id) from companies where deleted_at is null
	`

	companyDetailQuery = `
		select
			c.id,
			c.code,
			c.name,
			c.description,
			c.phone,
			c.email,
			c.website,
			c.address,
			ST_X(c.location::geometry) AS longitude,
			ST_Y(c.location::geometry) AS latitude,
			cty.name as country,
			lang.name as language,
			cry.name as currency,
			c.created_at,
			c.updated_at
		from companies c
		inner join countries cty on c.country_id = cty.id and cty.deleted_at is null 
		inner join languages lang on c.language_id = lang.id and lang.deleted_at is null
		inner join currencies cry on c.currency_id = cry.id and cry.deleted_at is null
		where c.id = $1 and c.deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, companyId uint) (*readmodel.CompanyDetail, error) {
	data := &readmodel.CompanyDetail{}

	if err := pg.QueryRowContext(ctx, companyDetailQuery, companyId).Scan(
		&data.Id,
		&data.Code,
		&data.Name,
		&data.Description,
		&data.Phone,
		&data.Email,
		&data.Website,
		&data.Address,
		&data.Location.Longitude,
		&data.Location.Latitude,
		&data.Country,
		&data.Language,
		&data.Currency,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCompanyNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.CompanyListResponse, error) {
	data := make([]*readmodel.CompanyList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, countCompanyListQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCompanyNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, companyListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		company := &readmodel.CompanyList{}

		company.Location = &readmodel.Point{}

		if err := rows.Scan(
			&company.Id,
			&company.Code,
			&company.Name,
			&company.Phone,
			&company.Email,
			&company.Website,
			&company.Location.Longitude,
			&company.Location.Latitude,
			&company.Country,
			&company.Language,
			&company.Currency,
			&company.CreatedAt,
			&company.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, company)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.CompanyListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Create(tx *gorm.DB, Company *entity.Company) error {
	result := tx.Create(Company)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrForeignKeyError
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.ErrDuplicateError
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, Company *entity.Company) error {
	result := tx.Model(&entity.Company{}).Updates(Company)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrForeignKeyError
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.ErrDuplicateError
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrCompanyNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, companyId uint) error {
	result := tx.Model(&entity.Company{}).Delete(companyId)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrCompanyNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
