package service

import (
	"context"
	"database/sql"

	"enterprise_resource_planning/services/material_management/dto"
	"enterprise_resource_planning/services/material_management/entity"
	"enterprise_resource_planning/services/material_management/readmodel"

	"gorm.io/gorm"
)

type AuditRepository interface {
	SetUserId(tx *gorm.DB, userId uint) error
}

type BrandRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.BrandListResponse, error)
	Detail(ctx context.Context, pg *sql.DB, brandId uint) (*readmodel.BrandDetail, error)
	Create(tx *gorm.DB, mg *entity.Brand) error
	Update(tx *gorm.DB, mg *entity.Brand) error
	Delete(tx *gorm.DB, brandId uint) error
}

type Service struct {
	pg              *gorm.DB
	auditRepository AuditRepository
	brandRepository BrandRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.BrandListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.brandRepository.List(ctx, pg, page)
}

func (s *Service) Detail(ctx context.Context, brandId uint) (*readmodel.BrandDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.brandRepository.Detail(ctx, pg, brandId)
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.BrandReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.brandRepository.Create(tx, &entity.Brand{
			Name: req.Name,
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, brandId uint, req *dto.BrandReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.brandRepository.Update(tx, &entity.Brand{
			Model: gorm.Model{
				ID: brandId,
			},
			Name: req.Name,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, brandId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.brandRepository.Delete(tx, brandId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, brandRepository BrandRepository) *Service {
	return &Service{
		pg:              pg,
		auditRepository: auditRepository,
		brandRepository: brandRepository,
	}
}
