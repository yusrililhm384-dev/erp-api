package repository

import (
	"errors"

	"enterprise_resource_planning/services/hr/entity"

	errs "enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

type Repository struct {
}

func (r *Repository) Create(tx *gorm.DB, mutaion *entity.EmployeeMutation) error {
	result := tx.Create(mutaion)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrForeignKeyError
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, mutaion *entity.EmployeeMutation) error {
	query := tx.Model(&entity.EmployeeMutation{})

	if mutaion.DeletedAt.Valid {
		result := query.Unscoped().Where("employee_id = ? and deleted_at is null", mutaion.EmployeeID).Updates(map[string]any{
			"end_date":   mutaion.EndDate,
			"updated_at": mutaion.UpdatedAt,
			"deleted_at": mutaion.DeletedAt.Time,
		})

		if err := result.Error; err != nil {
			if errors.Is(err, gorm.ErrForeignKeyViolated) {
				return errs.ErrForeignKeyError
			}

			return err
		}

		if result.RowsAffected == 0 {
			return errs.ErrMutationError
		}

		return nil
	}

	result := query.Where("employee_id = ? and deleted_at is null", mutaion.EmployeeID).Updates(mutaion)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrForeignKeyError
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrMutationError
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
