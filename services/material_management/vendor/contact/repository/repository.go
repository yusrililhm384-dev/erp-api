package repository

import (
	"context"
	"database/sql"
	"errors"

	"enterprise_resource_planning/services/material_management/entity"
	"enterprise_resource_planning/services/material_management/readmodel"

	errs "enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

const (
	vendorContactDetailQuery = `
		select id, name, email, phone, created_at, updated_at from vendor_contacts where id = $1 and vendor_id = $2 and deleted_at is null 
	`
)

type Repository struct {
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, contactId, vendorId uint) (*readmodel.VendorContactDetail, error) {
	data := &readmodel.VendorContactDetail{}

	if err := pg.QueryRowContext(ctx, vendorContactDetailQuery, contactId, vendorId).Scan(
		&data.Id,
		&data.Name,
		&data.Email,
		&data.Phone,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrVendorContactNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) Create(tx *gorm.DB, contact *entity.VendorContact) error {
	if err := tx.Create(contact).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrForeignKeyViolated):
			return errs.ErrForeignKeyError
		case errors.Is(err, gorm.ErrDuplicatedKey):
			return errs.ErrDuplicateVendorContact
		default:
			return err
		}
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, contact *entity.VendorContact) error {
	res := tx.Model(&entity.VendorContact{}).Where("id = ? and vendor_id = ?", contact.ID, contact.VendorID).Updates(contact)

	if err := res.Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrForeignKeyViolated):
			return errs.ErrForeignKeyError
		case errors.Is(err, gorm.ErrDuplicatedKey):
			return errs.ErrDuplicateVendorContact
		default:
			return err
		}
	}

	if res.RowsAffected == 0 {
		return errs.ErrVendorContactNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, contactId, vendorId uint) error {
	res := tx.Where("id = ? and vendor_id = ?", contactId, vendorId).Delete(&entity.VendorContact{})

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrVendorContactNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
