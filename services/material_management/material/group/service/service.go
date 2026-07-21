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

type MaterialGroupRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.MaterialGroupListResponse, error)
	Detail(ctx context.Context, pg *sql.DB, materialGroupId uint) (*readmodel.MaterialGroupDetail, error)
	Create(tx *gorm.DB, mg *entity.MaterialGroup) error
	Update(tx *gorm.DB, mg *entity.MaterialGroup) error
	Delete(tx *gorm.DB, materialGroupId uint) error
}

type Service struct {
	pg                      *gorm.DB
	auditRepository         AuditRepository
	materialGroupRepository MaterialGroupRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.MaterialGroupListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.materialGroupRepository.List(ctx, pg, page)
}

func (s *Service) Detail(ctx context.Context, materialGroupId uint) (*readmodel.MaterialGroupDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.materialGroupRepository.Detail(ctx, pg, materialGroupId)
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreateMaterialGroupReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.materialGroupRepository.Create(tx, &entity.MaterialGroup{
			Name:        req.Name,
			Description: req.Description,
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, materialGroupId uint, req *dto.UpdateMaterialGroupReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.materialGroupRepository.Update(tx, &entity.MaterialGroup{
			Model: gorm.Model{
				ID: materialGroupId,
			},
			Name:        req.Name,
			Description: req.Description,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, materialGroupId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.materialGroupRepository.Delete(tx, materialGroupId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, materialGroupRepository MaterialGroupRepository) *Service {
	return &Service{
		pg:                      pg,
		auditRepository:         auditRepository,
		materialGroupRepository: materialGroupRepository,
	}
}
