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
	materialTypeDetailQuery = `
		select id, name, description, created_at, updated_at from material_types where id = $1 and deleted_at is null
	`

	materialTypeListQuery = `
		select id, name, created_at, updated_at from material_types where deleted_at is null order by created_at desc limit $1 offset $2
	`

	countMaterialTypeListQuery = `
		select count(id) from material_types where deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.MaterialTypeListResponse, error) {
	data := make([]*readmodel.MaterialTypeList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, countMaterialTypeListQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrMaterialTypeNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, materialTypeListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		mt := &readmodel.MaterialTypeList{}

		if err := rows.Scan(
			&mt.Id,
			&mt.Name,
			&mt.CreatedAt,
			&mt.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, mt)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.MaterialTypeListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, materialTypeId uint) (*readmodel.MaterialTypeDetail, error) {
	data := &readmodel.MaterialTypeDetail{}

	if err := pg.QueryRowContext(ctx, materialTypeDetailQuery, materialTypeId).Scan(
		&data.Id,
		&data.Name,
		&data.Description,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrMaterialTypeNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) Create(tx *gorm.DB, mt *entity.MaterialType) error {
	if err := tx.Create(mt).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, mt *entity.MaterialType) error {
	res := tx.Model(&entity.MaterialType{}).Updates(mt)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrMaterialTypeNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, materialTypeId uint) error {
	res := tx.Model(&entity.MaterialType{}).Delete(materialTypeId)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrMaterialTypeNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
