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

type CurrencyRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.CurrencyListResponse, error)
	Detail(ctx context.Context, pg *sql.DB, currencyId uint) (*readmodel.CurrencyDetail, error)
	Create(tx *gorm.DB, c *entity.Currency) error
	Update(tx *gorm.DB, c *entity.Currency) error
	Delete(tx *gorm.DB, currencyId uint) error
}

type Service struct {
	pg                 *gorm.DB
	auditRepository    AuditRepository
	currencyRepository CurrencyRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.CurrencyListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.currencyRepository.List(ctx, pg, page)
}

func (s *Service) Detail(ctx context.Context, currencyId uint) (*readmodel.CurrencyDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.currencyRepository.Detail(ctx, pg, currencyId)
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreateCurrencyReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.currencyRepository.Create(tx, &entity.Currency{
			Code:      req.Code,
			Name:      req.Name,
			CountryID: req.CountryId,
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, currencyId uint, req *dto.UpdateCurrencyReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.currencyRepository.Update(tx, &entity.Currency{
			Model: gorm.Model{
				ID: currencyId,
			},
			Code:      req.Code,
			Name:      req.Name,
			CountryID: req.CountryId,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, currencyId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.currencyRepository.Delete(tx, currencyId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, currencyRepository CurrencyRepository) *Service {
	return &Service{
		pg:                 pg,
		auditRepository:    auditRepository,
		currencyRepository: currencyRepository,
	}
}
