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
	paymentMethodDetailQuery = `
		select 
			id, 
			code, 
			name,
			description,
			created_at, 
			updated_at 
		from 
			payment_method
		where id = $1 and deleted_at is null
	`

	paymentMethodListQuery = `
		select id, code, name, created_at, updated_at from payment_method where deleted_at is null order by created_at desc limit $1 offset $2
	`

	countPaymentMethodListQuery = `
		select count(id) from payment_method where deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.PaymentTermListResponse, error) {
	data := make([]*readmodel.PaymentTermList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, countPaymentMethodListQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPaymentMethodNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, paymentMethodListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		pm := &readmodel.PaymentTermList{}

		if err := rows.Scan(
			&pm.Id,
			&pm.Code,
			&pm.Name,
			&pm.CreatedAt,
			&pm.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, pm)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.PaymentTermListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, paymentMethodId uint) (*readmodel.PaymentTermDetail, error) {
	data := &readmodel.PaymentTermDetail{}

	if err := pg.QueryRowContext(ctx, paymentMethodDetailQuery, paymentMethodId).Scan(
		&data.Id,
		&data.Code,
		&data.Name,
		&data.Description,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPaymentMethodNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) Create(tx *gorm.DB, pm *entity.PaymentMethod) error {
	if err := tx.Create(pm).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrForeignKeyError
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.ErrDuplicateError
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, pm *entity.PaymentMethod) error {
	res := tx.Model(&entity.PaymentMethod{}).Updates(pm)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrForeignKeyError
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.ErrDuplicateError
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrPaymentMethodNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, paymentMethodId uint) error {
	res := tx.Model(&entity.PaymentMethod{}).Delete(paymentMethodId)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrPaymentMethodNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
