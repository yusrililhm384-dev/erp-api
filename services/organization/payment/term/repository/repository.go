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
	paymentTermDetailQuery = `
		select 
			id, 
			code, 
			name,
			description,
			created_at, 
			updated_at 
		from 
			payment_terms
		where id = $1 and deleted_at is null
	`

	paymentTermListQuery = `
		select id, code, name, created_at, updated_at from payment_terms where deleted_at is null order by created_at desc limit $1 offset $2
	`

	countPaymentTermListQuery = `
		select count(id) from payment_terms where deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.PaymentTermListResponse, error) {
	data := make([]*readmodel.PaymentTermList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, countPaymentTermListQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPaymentTermNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, paymentTermListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		pt := &readmodel.PaymentTermList{}

		if err := rows.Scan(
			&pt.Id,
			&pt.Code,
			&pt.Name,
			&pt.CreatedAt,
			&pt.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, pt)
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

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, paymentTermId uint) (*readmodel.PaymentTermDetail, error) {
	data := &readmodel.PaymentTermDetail{}

	if err := pg.QueryRowContext(ctx, paymentTermDetailQuery, paymentTermId).Scan(
		&data.Id,
		&data.Code,
		&data.Name,
		&data.Description,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPaymentTermNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) Create(tx *gorm.DB, pt *entity.PaymentTerm) error {
	if err := tx.Create(pt).Error; err != nil {
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

func (r *Repository) Update(tx *gorm.DB, pt *entity.PaymentTerm) error {
	res := tx.Model(&entity.PaymentTerm{}).Updates(pt)

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
		return errs.ErrPaymentTermNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, paymentTermId uint) error {
	res := tx.Model(&entity.PaymentTerm{}).Delete(paymentTermId)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrPaymentTermNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
