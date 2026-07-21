package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math"

	"enterprise_resource_planning/services/organization/entity"
	"enterprise_resource_planning/services/organization/readmodel"

	errs "enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

const (
	currencyDetailQuery = `
		select 
			cr.id, 
			cr.code, 
			cr.name,
			coalesce(
				jsonb_agg(
					jsonb_build_object(
						'id', cty.id,
						'code', cty.code,
						'name', cty.name
					)
				) filter (where cty.id is not null), '[]'::jsonb
			) as countries,
			cr.created_at, 
			cr.updated_at 
		from 
			currencies cr left join countries cty on cr.Currency_id = cty.id and cty.deleted_at is null
		where id = $1 and deleted_at is null
		group by cr.id
	`

	currencyListQuery = `
		select id, code, name, created_at, updated_at from currencies where deleted_at is null order by created_at desc limit $1 offset $2
	`

	countCurrencyListQuery = `
		select count(id) from currencies where deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.CurrencyListResponse, error) {
	data := make([]*readmodel.CurrencyList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, countCurrencyListQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCurrencyNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, currencyListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		c := &readmodel.CurrencyList{}

		if err := rows.Scan(
			&c.Id,
			&c.Code,
			&c.Name,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, c)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.CurrencyListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, currencyId uint) (*readmodel.CurrencyDetail, error) {
	data := &readmodel.CurrencyDetail{}

	var countries []byte

	if err := pg.QueryRowContext(ctx, currencyDetailQuery, currencyId).Scan(
		&data.Id,
		&data.Code,
		&data.Name,
		&countries,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCurrencyNotFound
		}

		return nil, err
	}

	if err := json.Unmarshal(countries, &data.Countries); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *Repository) Create(tx *gorm.DB, c *entity.Currency) error {
	if err := tx.Create(c).Error; err != nil {
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

func (r *Repository) Update(tx *gorm.DB, c *entity.Currency) error {
	res := tx.Model(&entity.Currency{}).Updates(c)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrForeignKeyError
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.ErrDuplicateError
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrCurrencyNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, currencyId uint) error {
	res := tx.Model(&entity.Currency{}).Delete(currencyId)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrCurrencyNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
