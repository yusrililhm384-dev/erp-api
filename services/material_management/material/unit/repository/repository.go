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
	materialUnitDetailQuery = `
		select mu.id, mu.code, mu.name, muc.name as category, mu.description, mu.created_at, mu.updated_at from materail_units mu left join material_unit_categories muc on mu.material_category_id = muc.id where mu.id = $1 and mu.deleted_at is null
	`

	materialUnitListQuery = `
		select mu.id, mu.code, mu.name, muc.name as category, mu.created_at, mu.updated_at from materail_units mu left join material_unit_categories muc on mu.material_category_id = muc.id where mu.deleted_at is null order by mu.created_at desc limit $1 offset $2
	`

	countMaterialUnitListQuery = `
		select count(id) from materail_units where deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.MaterialUnitListResponse, error) {
	data := make([]*readmodel.MaterialUnitList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, countMaterialUnitListQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrMaterialUnitNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, materialUnitListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		mu := &readmodel.MaterialUnitList{}

		if err := rows.Scan(
			&mu.Id,
			&mu.Code,
			&mu.Name,
			&mu.Category,
			&mu.CreatedAt,
			&mu.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, mu)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.MaterialUnitListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, materialUnitId uint) (*readmodel.MaterialUnitDetail, error) {
	data := &readmodel.MaterialUnitDetail{}

	if err := pg.QueryRowContext(ctx, materialUnitDetailQuery, materialUnitId).Scan(
		&data.Id,
		&data.Code,
		&data.Name,
		&data.Category,
		&data.Description,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrMaterialUnitNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) Create(tx *gorm.DB, mu *entity.MaterialUnit) error {
	if err := tx.Create(mu).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrForeignKeyError
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, mu *entity.MaterialUnit) error {
	res := tx.Model(&entity.MaterialUnit{}).Updates(mu)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrForeignKeyError
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrMaterialUnitNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, materialUnitId uint) error {
	res := tx.Model(&entity.MaterialUnit{}).Delete(materialUnitId)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrMaterialUnitNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
