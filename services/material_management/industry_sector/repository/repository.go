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
	industrySectorDetailQuery = `
		select id, name, description, created_at, updated_at from industry_sectors where id = $1 and deleted_at is null
	`

	industrySectorListQuery = `
		select id, name, created_at, updated_at from industry_sectors where deleted_at is null order by created_at desc limit $1 offset $2
	`

	countIndustrySectorListQuery = `
		select count(id) from industry_sectors where deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.IndustrySectorListResponse, error) {
	data := make([]*readmodel.IndustrySectorList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, countIndustrySectorListQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrIndustrySectorNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, industrySectorListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		is := &readmodel.IndustrySectorList{}

		if err := rows.Scan(
			&is.Id,
			&is.Name,
			&is.CreatedAt,
			&is.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, is)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.IndustrySectorListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, industrySectorId uint) (*readmodel.IndustrySectorDetail, error) {
	data := &readmodel.IndustrySectorDetail{}

	if err := pg.QueryRowContext(ctx, industrySectorDetailQuery, industrySectorId).Scan(
		&data.Id,
		&data.Name,
		&data.Description,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrIndustrySectorNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) Create(tx *gorm.DB, is *entity.IndustrySector) error {
	if err := tx.Create(is).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, is *entity.IndustrySector) error {
	res := tx.Model(&entity.IndustrySector{}).Updates(is)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrIndustrySectorNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, industrySectorId uint) error {
	res := tx.Model(&entity.IndustrySector{}).Delete(industrySectorId)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrIndustrySectorNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
