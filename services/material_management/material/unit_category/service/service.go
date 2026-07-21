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

type MaterialUnitCategoryRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.MaterialUnitCategoryListResponse, error)
	Detail(ctx context.Context, pg *sql.DB, materialUnitCategoryId uint) (*readmodel.MaterialUnitCategoryDetail, error)
	Create(tx *gorm.DB, bu *entity.MaterialUnitCategory) error
	Update(tx *gorm.DB, bu *entity.MaterialUnitCategory) error
	Delete(tx *gorm.DB, materialUnitCategoryId uint) error
}

type Service struct {
	pg                             *gorm.DB
	auditRepository                AuditRepository
	materialUnitCategoryRepository MaterialUnitCategoryRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.MaterialUnitCategoryListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.materialUnitCategoryRepository.List(ctx, pg, page)
}

func (s *Service) Detail(ctx context.Context, materialUnitCategoryId uint) (*readmodel.MaterialUnitCategoryDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.materialUnitCategoryRepository.Detail(ctx, pg, materialUnitCategoryId)
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreateMaterialUnitCategoryReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		muc := &entity.MaterialUnitCategory{
			Code: req.Code,
			Name: req.Name,
		}

		for _, mu := range req.MaterialUnits {
			muc.MaterialUnits = append(muc.MaterialUnits, &entity.MaterialUnit{
				Code:        mu.Code,
				Name:        mu.Name,
				Description: mu.Description,
			})
		}

		return s.materialUnitCategoryRepository.Create(tx, muc)
	})
}

func (s *Service) Update(ctx context.Context, userId uint, materialUnitCategoryId uint, req *dto.UpdateMaterialUnitCategoryReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.materialUnitCategoryRepository.Update(tx, &entity.MaterialUnitCategory{
			Model: gorm.Model{
				ID: materialUnitCategoryId,
			},
			Code: req.Code,
			Name: req.Name,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, materialUnitCategoryId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.materialUnitCategoryRepository.Delete(tx, materialUnitCategoryId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, materialUnitCategoryRepository MaterialUnitCategoryRepository) *Service {
	return &Service{
		pg:                             pg,
		auditRepository:                auditRepository,
		materialUnitCategoryRepository: materialUnitCategoryRepository,
	}
}
