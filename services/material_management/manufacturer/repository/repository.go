package repository

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"enterprise_resource_planning/services/material_management/entity"
	"enterprise_resource_planning/services/material_management/readmodel"

	errs "enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

const (
	manufacturerDetailQuery = `
		select id, code, name, created_at, updated_at from manufacturers where id = $1 and deleted_at is null
	`

	manufacturerListQuery = `
		select id, code, name, created_at, updated_at from manufacturers where deleted_at is null order by created_at desc limit $1 offset $2
	`

	countManufacturerListQuery = `
		select count(id) from manufacturers where deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.ManufacturerListResponse, error) {
	data := make([]*readmodel.ManufacturerList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, countManufacturerListQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrManufacturerNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, manufacturerListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		mf := &readmodel.ManufacturerList{}

		if err := rows.Scan(
			&mf.Id,
			&mf.Code,
			&mf.Name,
			&mf.CreatedAt,
			&mf.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, mf)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.ManufacturerListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, manufacturerId uint) (*readmodel.ManufacturerDetail, error) {
	data := &readmodel.ManufacturerDetail{}

	if err := pg.QueryRowContext(ctx, manufacturerDetailQuery, manufacturerId).Scan(
		&data.Id,
		&data.Code,
		&data.Name,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrManufacturerNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) Create(tx *gorm.DB, mf *entity.Manufacturer) error {
	if err := tx.Create(mf).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.ErrDuplicateError
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, mf *entity.Manufacturer) error {
	res := tx.Model(&entity.Manufacturer{}).Updates(mf)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.ErrDuplicateError
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrManufacturerNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, manufacturerId uint) error {
	res := tx.Model(&entity.Manufacturer{}).Delete(manufacturerId)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrManufacturerNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
