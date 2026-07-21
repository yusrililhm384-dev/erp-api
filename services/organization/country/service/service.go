package service

import (
	"context"
	"database/sql"

	"enterprise_resource_planning/services/organization/dto"
	"enterprise_resource_planning/services/organization/entity"
	"enterprise_resource_planning/services/organization/readmodel"

	"gorm.io/gorm"
)

type AuditRepository interface {
	SetUserId(tx *gorm.DB, userId uint) error
}

type CountryRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.CountryListResponse, error)
	Detail(ctx context.Context, pg *sql.DB, countryId uint) (*readmodel.CountryDetail, error)
	Create(tx *gorm.DB, c *entity.Country) error
	Update(tx *gorm.DB, c *entity.Country) error
	Delete(tx *gorm.DB, countryId uint) error
}

type Service struct {
	pg                *gorm.DB
	auditRepository   AuditRepository
	countryRepository CountryRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.CountryListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.countryRepository.List(ctx, pg, page)
}

func (s *Service) Detail(ctx context.Context, countryId uint) (*readmodel.CountryDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.countryRepository.Detail(ctx, pg, countryId)
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreateCountryReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.countryRepository.Create(tx, &entity.Country{
			Code: req.Code,
			Name: req.Name,
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, countryId uint, req *dto.UpdateCountryReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.countryRepository.Update(tx, &entity.Country{
			Model: gorm.Model{
				ID: countryId,
			},
			Code: req.Code,
			Name: req.Name,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, countryId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.countryRepository.Delete(tx, countryId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, countryRepository CountryRepository) *Service {
	return &Service{
		pg:                pg,
		auditRepository:   auditRepository,
		countryRepository: countryRepository,
	}
}
