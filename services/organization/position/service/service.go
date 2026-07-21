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

type PositionRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.PositionListResponse, error)
	Detail(ctx context.Context, pg *sql.DB, positionId uint) (*readmodel.PositionDetail, error)
	Create(tx *gorm.DB, mu *entity.Position) error
	Update(tx *gorm.DB, mu *entity.Position) error
	Delete(tx *gorm.DB, positionId uint) error
}

type Service struct {
	pg                 *gorm.DB
	auditRepository    AuditRepository
	positionRepository PositionRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.PositionListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.positionRepository.List(ctx, pg, page)
}

func (s *Service) Detail(ctx context.Context, positionId uint) (*readmodel.PositionDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.positionRepository.Detail(ctx, pg, positionId)
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreatePositionReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.positionRepository.Create(tx, &entity.Position{
			Name:               req.Name,
			Description:        req.Description,
			ApprovalID:         req.ApprovalId,
			PositionCategoryID: req.CategoryId,
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, positionId uint, req *dto.UpdatePositionReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.positionRepository.Update(tx, &entity.Position{
			Model: gorm.Model{
				ID: positionId,
			},
			Name:               req.Name,
			Description:        req.Description,
			ApprovalID:         req.ApprovalId,
			PositionCategoryID: req.CategoryId,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, positionId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.positionRepository.Delete(tx, positionId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, positionRepository PositionRepository) *Service {
	return &Service{
		pg:                 pg,
		auditRepository:    auditRepository,
		positionRepository: positionRepository,
	}
}
