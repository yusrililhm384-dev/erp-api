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
	departmentsListQuery = `
		select
			id,
			name,
			description,
			created_at,
			updated_at
		from departments
		where deleted_at is null
		order by created_at desc
		limit $1 offset $2
	`

	departmentCountQuery = `
		select count(id) from departments where deleted_at is null
	`

	departmentDetailQuery = `
		select
			id,
			name,
			description,
			created_at,
			updated_at
		from departments
		where id = $1 and deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, departmentId uint) (*readmodel.DepartmentDetail, error) {
	data := &readmodel.DepartmentDetail{}

	if err := pg.QueryRowContext(ctx, departmentDetailQuery, departmentId).Scan(
		&data.Id,
		&data.Name,
		&data.Description,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrDepartmentNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.DepartmentListResponse, error) {
	data := make([]*readmodel.DepartmentList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, departmentCountQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrDepartmentNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, departmentsListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		dept := &readmodel.DepartmentList{}

		if err := rows.Scan(
			&dept.Id,
			&dept.Name,
			&dept.Description,
			&dept.CreatedAt,
			&dept.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, dept)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.DepartmentListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Create(tx *gorm.DB, department *entity.Department) error {
	result := tx.Create(department)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, department *entity.Department) error {
	result := tx.Model(&entity.Department{}).Updates(department)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrDepartmentNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, departmentId uint) error {
	result := tx.Model(&entity.Department{}).Delete(departmentId)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrDepartmentNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
