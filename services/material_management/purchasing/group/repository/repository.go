package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math"

	"enterprise_resource_planning/services/material_management/entity"
	"enterprise_resource_planning/services/material_management/readmodel"

	errs "enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

const (
	purchasingGroupListQuery = `
		select
			pg.id,
			pg.code,
			pg.name,
			po.name as purchasing_organization
			pg.created_at,
			pg.updated_at
		from
			purchasing_groups pg
		inner join 
			purchasing_organizations po on pg.purchasing_organization_id = po.id and po.deleted_at is null
		where
			pg.deleted_at is null
		order by pg.created_at desc
		limit $1 offset $2
	`

	countPurchasingGroupListQuery = `
		select count(id) from purchasing_groups where deleted_at is null
	`

	purchasingGroupDetailQuery = `
		select
			pg.id,
			pg.code,
			pg.name,
			pg.description,
			po.name as purchasing_organization,
			coalesce(
				jsonb_agg(
					jsonb_build_object()
				) filter (where pgm.id is not null), '[]'::jsonb
			) as list,
			pg.created_at,
			pg.updated_at
		from
			purchasing_groups pg
		where
			pg.id = $1 and pg.deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, purchasingGroupId uint) (*readmodel.PurchasingGroupDetail, error) {
	data := &readmodel.PurchasingGroupDetail{}

	var list []byte

	if err := pg.QueryRowContext(ctx, purchasingGroupDetailQuery, purchasingGroupId).Scan(
		&data.Id,
		&data.Code,
		&data.Name,
		&data.Description,
		&data.PurchasingOrganization,
		&list,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPurchasingGroupNotFound
		}

		return nil, err
	}

	if err := json.Unmarshal(list, &data.PurchasingGroupMembers); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.PurchasingGroupListResponse, error) {
	data := make([]*readmodel.PurchasingGroupList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, countPurchasingGroupListQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPurchasingGroupNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, purchasingGroupListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		purchasingGroup := &readmodel.PurchasingGroupList{}

		if err := rows.Scan(
			&purchasingGroup.Id,
			&purchasingGroup.Code,
			&purchasingGroup.Name,
			&purchasingGroup.PurchasingOrganization,
			&purchasingGroup.CreatedAt,
			&purchasingGroup.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, purchasingGroup)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.PurchasingGroupListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Create(tx *gorm.DB, pg *entity.PurchasingGroup) error {
	result := tx.Create(pg)

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

func (r *Repository) Update(tx *gorm.DB, pg *entity.PurchasingGroup) error {
	result := tx.Model(&entity.PurchasingGroup{}).Updates(pg)

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
		return errs.ErrPurchasingGroupNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, purchasingGroupId uint) error {
	result := tx.Model(&entity.PurchasingGroup{}).Delete(purchasingGroupId)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrPurchasingGroupNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
