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
	materialGroupDetailQuery = `
		select id, name, description, created_at, updated_at from material_groups where id = $1 and deleted_at is null
	`

	materialGroupListQuery = `
		select id, name, created_at, updated_at from material_groups where deleted_at is null order by created_at desc limit $1 offset $2
	`

	countMaterialGroupListQuery = `
		select count(id) from material_groups where deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.MaterialGroupListResponse, error) {
	data := make([]*readmodel.MaterialGroupList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, countMaterialGroupListQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrMaterialGroupNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, materialGroupListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		mt := &readmodel.MaterialGroupList{}

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

	return &readmodel.MaterialGroupListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, materialGroupId uint) (*readmodel.MaterialGroupDetail, error) {
	data := &readmodel.MaterialGroupDetail{}

	if err := pg.QueryRowContext(ctx, materialGroupDetailQuery, materialGroupId).Scan(
		&data.Id,
		&data.Name,
		&data.Description,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrMaterialGroupNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) Create(tx *gorm.DB, mt *entity.MaterialGroup) error {
	if err := tx.Create(mt).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, mt *entity.MaterialGroup) error {
	res := tx.Model(&entity.MaterialGroup{}).Updates(mt)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrMaterialGroupNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, materialGroupId uint) error {
	res := tx.Model(&entity.MaterialGroup{}).Delete(materialGroupId)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrMaterialGroupNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
