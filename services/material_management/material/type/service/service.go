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

type MaterialTypeRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.MaterialTypeListResponse, error)
	Detail(ctx context.Context, pg *sql.DB, materialTypeId uint) (*readmodel.MaterialTypeDetail, error)
	Create(tx *gorm.DB, mt *entity.MaterialType) error
	Update(tx *gorm.DB, mt *entity.MaterialType) error
	Delete(tx *gorm.DB, materialTypeId uint) error
}

type Service struct {
	pg                     *gorm.DB
	auditRepository        AuditRepository
	materialTypeRepository MaterialTypeRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.MaterialTypeListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.materialTypeRepository.List(ctx, pg, page)
}

func (s *Service) Detail(ctx context.Context, materialTypeId uint) (*readmodel.MaterialTypeDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.materialTypeRepository.Detail(ctx, pg, materialTypeId)
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreateMaterialTypeReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.materialTypeRepository.Create(tx, &entity.MaterialType{
			Name:        req.Name,
			Description: req.Description,
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, materialTypeId uint, req *dto.UpdateMaterialTypeReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.materialTypeRepository.Update(tx, &entity.MaterialType{
			Model: gorm.Model{
				ID: materialTypeId,
			},
			Name:        req.Name,
			Description: req.Description,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, materialTypeId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.materialTypeRepository.Delete(tx, materialTypeId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, materialTypeRepository MaterialTypeRepository) *Service {
	return &Service{
		pg:                     pg,
		auditRepository:        auditRepository,
		materialTypeRepository: materialTypeRepository,
	}
}
