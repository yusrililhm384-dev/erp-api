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

type PositionCategoryRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.PositionCategoryListResponse, error)
	Detail(ctx context.Context, pg *sql.DB, positionCategoryId uint) (*readmodel.PositionCategoryDetail, error)
	Create(tx *gorm.DB, mu *entity.PositionCategory) error
	Update(tx *gorm.DB, mu *entity.PositionCategory) error
	Delete(tx *gorm.DB, positionCategoryId uint) error
}

type Service struct {
	pg                         *gorm.DB
	auditRepository            AuditRepository
	positionCategoryRepository PositionCategoryRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.PositionCategoryListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.positionCategoryRepository.List(ctx, pg, page)
}

func (s *Service) Detail(ctx context.Context, positionCategoryId uint) (*readmodel.PositionCategoryDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.positionCategoryRepository.Detail(ctx, pg, positionCategoryId)
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.PositionCategoryReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.positionCategoryRepository.Create(tx, &entity.PositionCategory{
			Name: req.Name,
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, positionCategoryId uint, req *dto.PositionCategoryReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.positionCategoryRepository.Update(tx, &entity.PositionCategory{
			Model: gorm.Model{
				ID: positionCategoryId,
			},
			Name: req.Name,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, positionCategoryId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.positionCategoryRepository.Delete(tx, positionCategoryId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, positionCategoryRepository PositionCategoryRepository) *Service {
	return &Service{
		pg:                         pg,
		auditRepository:            auditRepository,
		positionCategoryRepository: positionCategoryRepository,
	}
}
