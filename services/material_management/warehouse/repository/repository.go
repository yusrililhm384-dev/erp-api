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
	warehousesListQuery = `
		select
			id,
			name,
			description,
			address,
			ST_X(location::geometry) AS longitude,
			ST_Y(location::geometry) AS latitude,
			created_at,
			updated_at
		from warehouses
		where deleted_at is null
		order by created_at desc
		limit $1 offset $2
	`

	warehouseCountQuery = `
		select count(id) from warehouses where deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.WarehouseListResponse, error) {
	data := make([]*readmodel.WarehouseList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, warehouseCountQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrWarehouseNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, warehousesListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		warehouse := &readmodel.WarehouseList{}

		warehouse.Location = &readmodel.Point{}

		if err := rows.Scan(
			&warehouse.Id,
			&warehouse.Name,
			&warehouse.Description,
			&warehouse.Address,
			&warehouse.Location.Longitude,
			&warehouse.Location.Latitude,
			&warehouse.CreatedAt,
			&warehouse.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, warehouse)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.WarehouseListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Create(tx *gorm.DB, warehouse *entity.Warehouse) error {
	result := tx.Create(warehouse)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, warehouse *entity.Warehouse) error {
	result := tx.Model(&entity.Warehouse{}).Updates(warehouse)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrWarehouseNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, warehouseId uint) error {
	result := tx.Model(&entity.Warehouse{}).Delete(warehouseId)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrWarehouseNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
