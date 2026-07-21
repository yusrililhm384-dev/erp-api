package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math"

	"enterprise_resource_planning/services/material_management/entity"
	"enterprise_resource_planning/services/material_management/readmodel"

	errs "enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

const (
	materialCategoryUnitDetailQuery = `
		select 
			muc.id, 
			muc.code, 
			muc.name, 
			coalesce(
				jsonb_agg(
					jsonb_build_object(
						'code', mu.code,
						'name', mu.name
					)
				) filter (mu.id is not null), '[]'::jsonb
			) as list, 
			muc.created_at, 
			muc.updated_at 
		from materail_unit_categories muc left join material_unit muc on muc.material_category_id = muc.id 
		where muc.id = $1 and muc.deleted_at is null
		group by muc.id
	`

	materialUnitCategoryCategoryListQuery = `
		select id, code, name, created_at, updated_at from materail_unit_categories where deleted_at is null order by created_at desc limit $1 offset $2
	`

	countMaterialCategoryUnitListQuery = `
		select count(id) from materail_unit_categories where deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.MaterialUnitCategoryListResponse, error) {
	data := make([]*readmodel.MaterialUnitCategoryList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, countMaterialCategoryUnitListQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCategoryNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, materialUnitCategoryCategoryListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		muc := &readmodel.MaterialUnitCategoryList{}

		if err := rows.Scan(
			&muc.Id,
			&muc.Code,
			&muc.Name,
			&muc.CreatedAt,
			&muc.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, muc)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.MaterialUnitCategoryListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, materialCategoryUnitId uint) (*readmodel.MaterialUnitCategoryDetail, error) {
	data := &readmodel.MaterialUnitCategoryDetail{}

	var list []byte

	if err := pg.QueryRowContext(ctx, materialCategoryUnitDetailQuery, materialCategoryUnitId).Scan(
		&data.Id,
		&data.Code,
		&data.Name,
		&list,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCategoryNotFound
		}

		return nil, err
	}

	if err := json.Unmarshal(list, &data.List); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *Repository) Create(tx *gorm.DB, muc *entity.MaterialUnitCategory) error {
	if err := tx.Create(muc).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrForeignKeyError
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, muc *entity.MaterialUnitCategory) error {
	res := tx.Model(&entity.MaterialUnitCategory{}).Updates(muc)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrForeignKeyError
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrCategoryNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, materialCategoryUnitId uint) error {
	res := tx.Model(&entity.MaterialUnitCategory{}).Delete(materialCategoryUnitId)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrCategoryNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
