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
	brandDetailQuery = `
		select id, name, created_at, updated_at from brands where id = $1 and deleted_at is null
	`

	brandListQuery = `
		select id, name, created_at, updated_at from brands where deleted_at is null order by created_at desc limit $1 offset $2
	`

	countBrandListQuery = `
		select count(id) from brands where deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.BrandListResponse, error) {
	data := make([]*readmodel.BrandList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, countBrandListQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrBrandNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, brandListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		brand := &readmodel.BrandList{}

		if err := rows.Scan(
			&brand.Id,
			&brand.Name,
			&brand.CreatedAt,
			&brand.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, brand)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.BrandListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, brandId uint) (*readmodel.BrandDetail, error) {
	data := &readmodel.BrandDetail{}

	if err := pg.QueryRowContext(ctx, brandDetailQuery, brandId).Scan(
		&data.Id,
		&data.Name,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrBrandNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) Create(tx *gorm.DB, brand *entity.Brand) error {
	if err := tx.Create(brand).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, brand *entity.Brand) error {
	res := tx.Model(&entity.Brand{}).Updates(brand)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrBrandNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, brandId uint) error {
	res := tx.Model(&entity.Brand{}).Delete(brandId)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrBrandNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
