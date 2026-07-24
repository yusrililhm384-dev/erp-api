package repository

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"enterprise_resource_planning/services/material_management/entity"
	"enterprise_resource_planning/services/material_management/readmodel"

	errs "enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

const (
	purchasingOrganizationListQuery = `
		select
			po.id,
			po.code,
			po.name,
			c.name as company,
			cty.name as country,
			l.name as language,
			cry.name as currency,
			po.created_at,
			po.updated_at
		from
			purchasing_organizations po
		inner join 
			companies c on po.company_id = c.id and c.deleted_at is null
		inner join
			countries cty on po.country_id = cty.id and cty.deleted_at is null
		inner join
			languages l on po.language_id = l.id and l.deleted_at is null
		inner join
			currencies cry on po.currency_id = cry.id and cry.deleted_at is null
		where
			po.deleted_at is null
		order by po.created_at desc
		limit $1 offset $2
	`

	countPurchasingOrganizationListQuery = `
		select count(id) from purchasing_organizations where deleted_at is null
	`

	purchasingOrganizationDetailQuery = `
		select
			po.id,
			po.code,
			po.name,
			c.name as company,
			cty.name as country,
			l.name as language,
			cry.name as currency,
			po.created_at,
			po.updated_at
		from
			purchasing_organizations po
		inner join 
			companies c on po.company_id = c.id and c.deleted_at is null
		inner join
			countries cty on po.country_id = cty.id and cty.deleted_at is null
		inner join
			languages l on po.language_id = l.id and l.deleted_at is null
		inner join
			currencies cry on po.currency_id = cry.id and cry.deleted_at is null
		where
			po.id = $1 and po.deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, purchasingOrganizationId uint) (*readmodel.PurchasingOrganizationDetail, error) {
	data := &readmodel.PurchasingOrganizationDetail{}

	if err := pg.QueryRowContext(ctx, purchasingOrganizationDetailQuery, purchasingOrganizationId).Scan(
		&data.Id,
		&data.Code,
		&data.Name,
		&data.Description,
		&data.Company,
		&data.Country,
		&data.Language,
		&data.Currency,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPurchasingOrganizationNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.PurchasingOrganizationListResponse, error) {
	data := make([]*readmodel.PurchasingOrganizationList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, countPurchasingOrganizationListQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPurchasingOrganizationNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, purchasingOrganizationListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		po := &readmodel.PurchasingOrganizationList{}

		if err := rows.Scan(
			&po.Id,
			&po.Code,
			&po.Name,
			&po.Company,
			&po.Country,
			&po.Language,
			&po.Currency,
			&po.CreatedAt,
			&po.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, po)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.PurchasingOrganizationListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Create(tx *gorm.DB, po *entity.PurchasingOrganization) error {
	result := tx.Create(po)

	if err := result.Error; err != nil {
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

func (r *Repository) Update(tx *gorm.DB, po *entity.PurchasingOrganization) error {
	result := tx.Model(&entity.PurchasingOrganization{}).Updates(po)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrForeignKeyError
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.ErrDuplicateError
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrPurchasingOrganizationNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, purchasingOrganizationId uint) error {
	result := tx.Model(&entity.PurchasingOrganization{}).Delete(purchasingOrganizationId)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrPurchasingOrganizationNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
