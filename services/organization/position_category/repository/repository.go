package repository

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"enterprise_resource_planning/services/organization/entity"
	"enterprise_resource_planning/services/organization/readmodel"

	errs "enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

const (
	positionCategoriesListQuery = `
		select
			id,
			name,
			created_at,
			updated_at
		from position_categories
		where deleted_at is null
		order by created_at desc
		limit $1 offset $2
	`

	positionCategoriesCountQuery = `
		select count(id) from position_categories where deleted_at is null
	`

	positionCategoriesDetailQuery = `
		select
			id,
			name,
			created_at,
			updated_at
		from position_categories
		where id = $1 and deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, positionCategoryId uint) (*readmodel.PositionCategoryDetail, error) {
	data := &readmodel.PositionCategoryDetail{}

	if err := pg.QueryRowContext(ctx, positionCategoriesDetailQuery, positionCategoryId).Scan(
		&data.Id,
		&data.Name,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCategoryNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.PositionCategoryListResponse, error) {
	data := make([]*readmodel.PositionCategoryList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, positionCategoriesCountQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCategoryNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, positionCategoriesListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		pc := &readmodel.PositionCategoryList{}

		if err := rows.Scan(
			&pc.Id,
			&pc.Name,
			&pc.CreatedAt,
			&pc.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, pc)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.PositionCategoryListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Create(tx *gorm.DB, pc *entity.PositionCategory) error {
	result := tx.Create(pc)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, pc *entity.PositionCategory) error {
	result := tx.Model(&entity.PositionCategory{}).Updates(pc)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrCategoryNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, positionCategoryId uint) error {
	result := tx.Model(&entity.PositionCategory{}).Delete(positionCategoryId)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrCategoryNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
