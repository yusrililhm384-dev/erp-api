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
	languageDetailQuery = `
		select id, code, name, created_at, updated_at from languages where id = $1 and deleted_at is null
	`

	languageListQuery = `
		select id, code, name, created_at, updated_at from languages where deleted_at is null order by created_at desc limit $1 offset $2
	`

	countLanguageListQuery = `
		select count(id) from languages where deleted_at is null
	`
)

type Repository struct {
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.CountryListResponse, error) {
	data := make([]*readmodel.CountryList, 0)

	const limit = 50

	offset := (page - 1) * limit

	var count uint

	if err := pg.QueryRowContext(ctx, countLanguageListQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrLanguageNotFound
		}

		return nil, err
	}

	rows, err := pg.QueryContext(ctx, languageListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		language := &readmodel.CountryList{}

		if err := rows.Scan(
			&language.Id,
			&language.Code,
			&language.Name,
			&language.CreatedAt,
			&language.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, language)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.CountryListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, languageId uint) (*readmodel.CountryDetail, error) {
	data := &readmodel.CountryDetail{}

	if err := pg.QueryRowContext(ctx, languageDetailQuery, languageId).Scan(
		&data.Id,
		&data.Code,
		&data.Name,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrLanguageNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) Create(tx *gorm.DB, language *entity.Language) error {
	if err := tx.Create(language).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.ErrDuplicateError
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, language *entity.Language) error {
	res := tx.Model(&entity.Language{}).Updates(language)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.ErrDuplicateError
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrLanguageNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, languageId uint) error {
	res := tx.Model(&entity.Language{}).Delete(languageId)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrLanguageNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
