package service

import (
	"context"
	"database/sql"
	"fmt"

	"enterprise_resource_planning/services/organization/dto"
	"enterprise_resource_planning/services/organization/entity"
	"enterprise_resource_planning/services/organization/readmodel"

	"gorm.io/gorm"
)

type AuditRepository interface {
	SetUserId(tx *gorm.DB, userId uint) error
}

type CompanyRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.CompanyListResponse, error)
	Detail(ctx context.Context, pg *sql.DB, companyId uint) (*readmodel.CompanyDetail, error)
	Create(tx *gorm.DB, company *entity.Company) error
	Update(tx *gorm.DB, company *entity.Company) error
	Delete(tx *gorm.DB, companyId uint) error
}

type Service struct {
	pg                *gorm.DB
	auditRepository   AuditRepository
	companyRepository CompanyRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.CompanyListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.companyRepository.List(ctx, pg, page)
}

func (s *Service) Detail(ctx context.Context, companyId uint) (*readmodel.CompanyDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.companyRepository.Detail(ctx, pg, companyId)
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreateCompanyReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.companyRepository.Create(tx, &entity.Company{
			Name:        req.Name,
			Description: req.Description,
			Address:     req.Address,
			Location:    fmt.Sprintf("SRID=4326;POINT(%f %f)", req.Location.Longitude, req.Location.Latitude),
			Code:        req.Code,
			Phone:       req.Phone,
			Email:       req.Email,
			Website:     req.Website,
			CountryID:   req.CountryId,
			CurrencyID:  req.CurrencyId,
			LanguageID:  req.LanguageId,
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, companyId uint, req *dto.UpdateCompanyReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.companyRepository.Update(tx, &entity.Company{
			Model: gorm.Model{
				ID: companyId,
			},
			Name:        req.Name,
			Description: req.Description,
			Address:     req.Address,
			Location:    fmt.Sprintf("SRID=4326;POINT(%f %f)", req.Location.Longitude, req.Location.Latitude),
			Code:        req.Code,
			Phone:       req.Phone,
			Email:       req.Email,
			Website:     req.Website,
			CountryID:   req.CountryId,
			CurrencyID:  req.CurrencyId,
			LanguageID:  req.LanguageId,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, companyId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.companyRepository.Delete(tx, companyId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, companyRepository CompanyRepository) *Service {
	return &Service{
		pg:                pg,
		auditRepository:   auditRepository,
		companyRepository: companyRepository,
	}
}
