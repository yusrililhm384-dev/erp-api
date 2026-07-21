package repository

import (
	"context"
	"database/sql"
	"errors"

	"enterprise_resource_planning/services/hr/entity"
	"enterprise_resource_planning/services/hr/readmodel"

	errs "enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

const (
	findTypeDetailQuery = `select id, name, description, created_at, updated_at from types where id = $1 and deleted_at is null`

	findTypeListQuery = `select id, name, created_at, updated_at from types where deleted_at is null`
)

type Repository struct {
}

func (r *Repository) Create(tx *gorm.DB, t *entity.EmployeeType) error {
	result := tx.Create(t)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, t *entity.EmployeeType) error {
	result := tx.Model(&entity.EmployeeType{}).Updates(t)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrTypeNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, typeId uint) error {
	result := tx.Model(&entity.EmployeeType{}).Delete(typeId)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrTypeNotFound
	}

	return nil
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, typeId uint) (*entity.EmployeeType, error) {
	data := &entity.EmployeeType{}

	if err := pg.QueryRowContext(ctx, findTypeDetailQuery, typeId).Scan(
		&data.ID,
		&data.Name,
		&data.Description,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrTypeNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) List(ctx context.Context, pg *sql.DB) ([]*readmodel.Type, error) {
	data := make([]*readmodel.Type, 0)

	rows, err := pg.QueryContext(ctx, findTypeListQuery)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		t := &readmodel.Type{}

		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, t)
	}

	return data, nil
}

func New() *Repository {
	return &Repository{}
}
