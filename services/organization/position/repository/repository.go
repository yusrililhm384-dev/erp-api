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
	positionsListQuery = `
		select
			p.id,
			p.name,
			pc.name as category,
			coalesce(ap.name, 'C-Level') as approved_by,
			p.created_at,
			p.updated_at
		from positions p
		left join position_categories pc on p.position_category_id = pc.id
		left join positions ap on p.approval_id = ap.id
		where p.deleted_at is null
		order by p.created_at desc
		limit $1 offset $2
	`

	positionDetailQuery = `
		select
			p.id,
			p.name,
			p.description,
			pc.name as category,
			coalesce(ap.name, 'C-Level') as approved_by,
			p.created_at,
			p.updated_at
		from positions p
		left join position_categories pc on p.position_category_id = pc.id
		left join positions ap on p.approval_id = ap.id
		where p.id = $1 and p.deleted_at is null
	`

	countPositionQuery = `
		select count(id) from positions where deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, positionId uint) (*readmodel.PositionDetail, error) {
	data := &readmodel.PositionDetail{}

	if err := pg.QueryRowContext(ctx, positionDetailQuery, positionId).Scan(
		&data.Id,
		&data.Name,
		&data.Description,
		&data.Category,
		&data.ApprovedBy,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPositionNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.PositionListResponse, error) {
	data := make([]*readmodel.PositionList, 0)

	var count uint

	if err := pg.QueryRowContext(ctx, countPositionQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPositionNotFound
		}

		return nil, err
	}

	const limit = 50

	offset := (page - 1) * limit

	rows, err := pg.QueryContext(ctx, positionsListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		list := &readmodel.PositionList{}

		if err := rows.Scan(
			&list.Id,
			&list.Name,
			&list.Category,
			&list.ApprovedBy,
			&list.CreatedAt,
			&list.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, list)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.PositionListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Create(tx *gorm.DB, position *entity.Position) error {
	if err := tx.Create(position).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrForeignKeyError
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, position *entity.Position) error {
	res := tx.Model(&entity.Position{}).Updates(position)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrForeignKeyError
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrPositionNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, positionId uint) error {
	res := tx.Model(&entity.Position{}).Delete(positionId)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrPositionNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
