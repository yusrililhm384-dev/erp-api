package repository

import (
	"context"
	"database/sql"
	"errors"

	"enterprise_resource_planning/services/hr/entity"

	errs "enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

const (
	findByUsernameQuery = `select employee_id as id, username, password from users where username = $1 and deleted_at is null`
)

type Repository struct {
}

func (r *Repository) Create(tx *gorm.DB, user *entity.User) error {
	result := tx.Create(user)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.ErrDuplicateUsernameError
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, user *entity.User) error {
	result := tx.Model(&entity.User{}).Where("employee_id = ?", user.EmployeeID).Updates(user)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.ErrDuplicateUsernameError
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}

func (r *Repository) FindByUsername(ctx context.Context, pg *sql.DB, username string) (*entity.User, error) {
	user := &entity.User{}

	if err := pg.QueryRowContext(ctx, findByUsernameQuery, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func New() *Repository {
	return &Repository{}
}
