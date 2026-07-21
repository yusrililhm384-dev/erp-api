package repository

import (
	"errors"

	"enterprise_resource_planning/services/hr/entity"

	errs "enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

type Repository struct {
}

func (r *Repository) Create(tx *gorm.DB, legal *entity.EmployeeLegal) error {
	result := tx.Create(legal)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.ErrDuplicateLegalError
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, legal *entity.EmployeeLegal) error {
	result := tx.Model(&entity.EmployeeLegal{}).Where("employee_id = ?", legal.EmployeeID).Updates(legal)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.ErrDuplicateLegalError
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrEmployeeMaserNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
