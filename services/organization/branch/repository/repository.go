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
	branchesListQuery = `
		select
			b.id,
			c.name as company,
			b.name,
			b.description,
			b.address,
			ST_X(b.location::geometry) AS longitude,
			ST_Y(b.location::geometry) AS latitude,
			b.created_at,
			b.updated_at
		from branches b inner join companies c on b.company_id = c.id and c.deleted_at is null 
		where b.deleted_at is null
		order by b.created_at desc
		limit $1 offset $2
	`

	branchCountQuery = `
		select count(id) from branches where deleted_at is null
	`

	branchDetailQuery = `
		select
			b.id,
			c.name as company,
			b.name,
			b.description,
			b.address,
			ST_X(b.location::geometry) AS longitude,
			ST_Y(b.location::geometry) AS latitude,
			b.created_at,
			b.updated_at
		from branches b inner join companies c on b.company_id = c.id and c.deleted_at is null
		where b.id = $1 and b.deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, branchId uint) (*readmodel.BranchDetail, error) {
	data := &readmodel.BranchDetail{}

	if err := pg.QueryRowContext(ctx, branchDetailQuery, branchId).Scan(
		&data.Id,
		&data.Company,
		&data.Name,
		&data.Description,
		&data.Address,
		&data.Location.Longitude,
		&data.Location.Latitude,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrBranchNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.BranchListResponse, error) {
	data := make([]*readmodel.BranchList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, branchCountQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrBranchNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, branchesListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		branch := &readmodel.BranchList{}

		branch.Location = &readmodel.Point{}

		if err := rows.Scan(
			&branch.Id,
			&branch.Company,
			&branch.Name,
			&branch.Description,
			&branch.Address,
			&branch.Location.Longitude,
			&branch.Location.Latitude,
			&branch.CreatedAt,
			&branch.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, branch)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.BranchListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Create(tx *gorm.DB, branch *entity.Branch) error {
	result := tx.Create(branch)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrForeignKeyError
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, branch *entity.Branch) error {
	result := tx.Model(&entity.Branch{}).Updates(branch)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrForeignKeyError
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrBranchNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, branchId uint) error {
	result := tx.Model(&entity.Branch{}).Delete(branchId)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrBranchNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
