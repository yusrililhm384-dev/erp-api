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
	storageLocationDetailQuery = `
		select id, code, name, description, created_at, updated_at from storage_locations where id = $1 and deleted_at is null
	`

	storageLocationListQuery = `
		select id code, name, created_at, updated_at from storage_locations where deleted_at is null order by created_at desc limit $1 offset $2
	`

	countStorageLocationListQuery = `
		select count(id) from storage_locations where deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.StorageLocationListResponse, error) {
	data := make([]*readmodel.StorageLocationList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, countStorageLocationListQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrStorageLocationNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, storageLocationListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		sl := &readmodel.StorageLocationList{}

		if err := rows.Scan(
			&sl.Id,
			&sl.Code,
			&sl.Name,
			&sl.CreatedAt,
			&sl.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, sl)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.StorageLocationListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, storageLocationId uint) (*readmodel.StorageLocationDetail, error) {
	data := &readmodel.StorageLocationDetail{}

	if err := pg.QueryRowContext(ctx, storageLocationDetailQuery, storageLocationId).Scan(
		&data.Id,
		&data.Code,
		&data.Name,
		&data.Description,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrStorageLocationNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) Create(tx *gorm.DB, sl *entity.StorageLocation) error {
	if err := tx.Create(sl).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrForeignKeyViolated):
			return errs.ErrWarehouseNotFound
		case errors.Is(err, gorm.ErrDuplicatedKey):
			return errs.ErrDuplicateStorageLocationCode
		default:
			return err
		}
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, sl *entity.StorageLocation) error {
	res := tx.Model(&entity.StorageLocation{}).Updates(sl)

	if err := res.Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrForeignKeyViolated):
			return errs.ErrUserNotFound
		case errors.Is(err, gorm.ErrDuplicatedKey):
			return errs.ErrDuplicateStorageLocationCode
		default:
			return err
		}
	}

	if res.RowsAffected == 0 {
		return errs.ErrStorageLocationNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, storageLocationId uint) error {
	res := tx.Model(&entity.StorageLocation{}).Delete(storageLocationId)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrStorageLocationNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
