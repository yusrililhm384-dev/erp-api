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
	materialStatusDetailQuery = `
		select id, code, name, created_at, updated_at from statuses where id = $1 and deleted_at is null
	`

	materialStatusListQuery = `
		select id, code, name, created_at, updated_at from statuses where deleted_at is null order by created_at desc limit $1 offset $2
	`

	countMaterialStatusListQuery = `
		select count(id) from statuses where deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.StatusListResponse, error) {
	data := make([]*readmodel.StatusList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, countMaterialStatusListQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrMaterialStatusNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, materialStatusListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		ms := &readmodel.StatusList{}

		if err := rows.Scan(
			&ms.Id,
			&ms.Code,
			&ms.Name,
			&ms.CreatedAt,
			&ms.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, ms)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.StatusListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, statusId uint) (*readmodel.StatusDetail, error) {
	data := &readmodel.StatusDetail{}

	if err := pg.QueryRowContext(ctx, materialStatusDetailQuery, statusId).Scan(
		&data.Id,
		&data.Code,
		&data.Name,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrMaterialStatusNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) Create(tx *gorm.DB, ms *entity.Status) error {
	if err := tx.Create(ms).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, ms *entity.Status) error {
	res := tx.Model(&entity.Status{}).Updates(ms)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrMaterialStatusNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, statusId uint) error {
	res := tx.Model(&entity.Status{}).Delete(statusId)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrMaterialStatusNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
