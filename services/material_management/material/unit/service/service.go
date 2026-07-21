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

type MaterialUnitRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.MaterialUnitListResponse, error)
	Detail(ctx context.Context, pg *sql.DB, materialUnitId uint) (*readmodel.MaterialUnitDetail, error)
	Create(tx *gorm.DB, bu *entity.MaterialUnit) error
	Update(tx *gorm.DB, bu *entity.MaterialUnit) error
	Delete(tx *gorm.DB, materialUnitId uint) error
}

type Service struct {
	pg                     *gorm.DB
	auditRepository        AuditRepository
	materailUnitRepository MaterialUnitRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.MaterialUnitListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.materailUnitRepository.List(ctx, pg, page)
}

func (s *Service) Detail(ctx context.Context, materialUnitId uint) (*readmodel.MaterialUnitDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.materailUnitRepository.Detail(ctx, pg, materialUnitId)
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreateMaterialUnitReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.materailUnitRepository.Create(tx, &entity.MaterialUnit{
			Code:                   req.Code,
			Name:                   req.Name,
			Description:            req.Description,
			MaterialUnitCategoryID: req.CategoryId,
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, materialUnitId uint, req *dto.UpdateMaterialUnitReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.materailUnitRepository.Update(tx, &entity.MaterialUnit{
			Model: gorm.Model{
				ID: materialUnitId,
			},
			Code:                   req.Code,
			Name:                   req.Name,
			Description:            req.Description,
			MaterialUnitCategoryID: req.CategoryId,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, materialUnitId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.materailUnitRepository.Delete(tx, materialUnitId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, materialUnitRepository MaterialUnitRepository) *Service {
	return &Service{
		pg:                     pg,
		auditRepository:        auditRepository,
		materailUnitRepository: materialUnitRepository,
	}
}
