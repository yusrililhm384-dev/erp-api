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

type IndustrySectorRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.IndustrySectorListResponse, error)
	Detail(ctx context.Context, pg *sql.DB, industrySectorId uint) (*readmodel.IndustrySectorDetail, error)
	Create(tx *gorm.DB, mg *entity.IndustrySector) error
	Update(tx *gorm.DB, mg *entity.IndustrySector) error
	Delete(tx *gorm.DB, industrySectorId uint) error
}

type Service struct {
	pg                       *gorm.DB
	auditRepository          AuditRepository
	industrySectorRepository IndustrySectorRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.IndustrySectorListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.industrySectorRepository.List(ctx, pg, page)
}

func (s *Service) Detail(ctx context.Context, industrySectorId uint) (*readmodel.IndustrySectorDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.industrySectorRepository.Detail(ctx, pg, industrySectorId)
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreateIndustrySectorReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.industrySectorRepository.Create(tx, &entity.IndustrySector{
			Name:        req.Name,
			Description: req.Description,
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, industrySectorId uint, req *dto.UpdateIndustrySectorReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.industrySectorRepository.Update(tx, &entity.IndustrySector{
			Model: gorm.Model{
				ID: industrySectorId,
			},
			Name:        req.Name,
			Description: req.Description,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, industrySectorId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.industrySectorRepository.Delete(tx, industrySectorId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, industrySectorRepository IndustrySectorRepository) *Service {
	return &Service{
		pg:                       pg,
		auditRepository:          auditRepository,
		industrySectorRepository: industrySectorRepository,
	}
}
