package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math"

	"enterprise_resource_planning/services/organization/entity"
	"enterprise_resource_planning/services/organization/readmodel"

	errs "enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

const (
	vendorListQuery = `
		select
			id,
			name,
			description,
			address,
			ST_X(location::geometry) AS longitude,
			ST_Y(location::geometry) AS latitude,
			created_at,
			updated_at
		from
			vendors
		where deleted_at is null
		order by created_at desc
		limit $1 offset $2
	`

	countVendorQuery = `
		select count(id) from vendors where deleted_at is null
	`

	vendorDetailQuery = `
		select
			v.id,
			v.name,
			v.description,
			v.address,
			ST_X(v.location::geometry) AS longitude,
			ST_Y(v.location::geometry) AS latitude,
			coalesce(
				jsonb_agg(
					jsonb_build_object(
						'id', c.id,
						'name', c.name,
						'email', c.email,
						'phone', c.phone
					)
				) filter (where c.id is not null), '[]'::jsonb
			) as contacts,
			v.created_at,
			v.updated_at
		from
			vendors v left join vendor_contacts c on c.vendor_id = v.id and c.deleted_at is null
		where v.id = $1 and v.deleted_at is null group by v.id
	`
)

type Repository struct {
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, vendorId uint) (*readmodel.VendorDetail, error) {
	data := &readmodel.VendorDetail{}

	data.Location = &readmodel.Point{}

	var contacts []byte

	if err := pg.QueryRowContext(ctx, vendorDetailQuery, vendorId).Scan(
		&data.Id,
		&data.Name,
		&data.Description,
		&data.Address,
		&data.Location.Longitude,
		&data.Location.Latitude,
		&contacts,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrVendorNotFound
		}

		return nil, err
	}

	if err := json.Unmarshal(contacts, &data.Contacts); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.VendorListResponse, error) {
	data := make([]*readmodel.VendorList, 0)

	var count uint

	if err := pg.QueryRowContext(ctx, countVendorQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrVendorNotFound
		}

		return nil, err
	}

	const limit = 50

	offset := (page - 1) * limit

	rows, err := pg.QueryContext(ctx, vendorListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		list := &readmodel.VendorList{}

		if err := rows.Scan(
			&list.Id,
			&list.Name,
			&list.Description,
			&list.Address,
			&list.Location,
			&list.CreatedAt,
			&list.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, list)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.VendorListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Create(tx *gorm.DB, vendor *entity.Vendor) error {
	if err := tx.Create(vendor).Error; err != nil {
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

func (r *Repository) Update(tx *gorm.DB, vendor *entity.Vendor) error {
	res := tx.Model(&entity.Vendor{}).Updates(vendor)

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
		return errs.ErrVendorNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, vendorId uint) error {
	res := tx.Model(&entity.Vendor{}).Delete(vendorId)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrVendorNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
