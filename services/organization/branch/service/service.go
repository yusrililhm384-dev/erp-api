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

type BranchRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.BranchListResponse, error)
	Detail(ctx context.Context, pg *sql.DB, branchId uint) (*readmodel.BranchDetail, error)
	Create(tx *gorm.DB, mu *entity.Branch) error
	Update(tx *gorm.DB, mu *entity.Branch) error
	Delete(tx *gorm.DB, branchId uint) error
}

type Service struct {
	pg              *gorm.DB
	auditRepository AuditRepository
	brachRepository BranchRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.BranchListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.brachRepository.List(ctx, pg, page)
}

func (s *Service) Detail(ctx context.Context, branchId uint) (*readmodel.BranchDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.brachRepository.Detail(ctx, pg, branchId)
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreateBranchReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.brachRepository.Create(tx, &entity.Branch{
			CompanyID:   req.CompanyId,
			Name:        req.Name,
			Description: req.Description,
			Address:     req.Address,
			Location:    fmt.Sprintf("SRID=4326;POINT(%f %f)", req.Location.Longitude, req.Location.Latitude),
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, branchId uint, req *dto.UpdateBranchReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.brachRepository.Update(tx, &entity.Branch{
			Model: gorm.Model{
				ID: branchId,
			},
			CompanyID:   req.CompanyId,
			Name:        req.Name,
			Description: req.Description,
			Address:     req.Address,
			Location:    fmt.Sprintf("SRID=4326;POINT(%f %f)", req.Location.Longitude, req.Location.Latitude),
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, branchId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.brachRepository.Delete(tx, branchId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, branchRepository BranchRepository) *Service {
	return &Service{
		pg:              pg,
		auditRepository: auditRepository,
		brachRepository: branchRepository,
	}
}
