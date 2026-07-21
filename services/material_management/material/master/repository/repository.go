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
	materialMasterListQuery = `
		select
			mm.id,
			mm.code,
			mm.name,
			mt.name as material_type,
			mg.name as material_group,
			mu.name as base_unit,
			ms.name as status,
			mm.created_at,
			mm.updated_at
		from
			material_masters mm
		inner join material_types mt on mm.material_type_id = mt.id and mt.deleted_at is null
		inner join material_groups mg on mm.material_group_id = mg.id and mg.deleted_at is null
		inner join material_units mu on mm.base_unit_id = mu.id and mu.deleted_at is null
		inner join statuses ms on mm.material_status_id = ms.id and ms.deleted_at is null
		where
			mm.deleted_at is null
		order by created_at desc
		limit $1 offset $2
	`

	materialMasterDetailQuery = `
		select
			
		from
		where
		group by
	`

	countMaterialMasterListQuery = `
		select count(id) from material_masters where deleted_at
	`
)

type Repository struct {
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.MaterialMasterListResponse, error) {
	data := make([]*readmodel.MaterialMasterList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, countMaterialMasterListQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrMaterialMasterNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, materialMasterListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		mm := &readmodel.MaterialMasterList{}

		if err := rows.Scan(
			&mm.Id,
			&mm.Code,
			&mm.Name,
			&mm.MaterialType,
			&mm.MaterialGroup,
			&mm.BaseUnit,
			&mm.Status,
			&mm.CreatedAt,
			&mm.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, mm)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.MaterialMasterListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Create(tx *gorm.DB, mm *entity.MaterialMaster) error {
	if err := tx.Create(mm).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrForeignKeyError
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, mm *entity.MaterialMaster) error {
	res := tx.Model(&entity.MaterialMaster{}).Updates(mm)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrForeignKeyError
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrMaterialMasterNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
