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

type StatusRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.StatusListResponse, error)
	Detail(ctx context.Context, pg *sql.DB, statusId uint) (*readmodel.StatusDetail, error)
	Create(tx *gorm.DB, ms *entity.Status) error
	Update(tx *gorm.DB, ms *entity.Status) error
	Delete(tx *gorm.DB, statusId uint) error
}

type Service struct {
	pg               *gorm.DB
	auditRepository  AuditRepository
	statusRepository StatusRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.StatusListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.statusRepository.List(ctx, pg, page)
}

func (s *Service) Detail(ctx context.Context, statusId uint) (*readmodel.StatusDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.statusRepository.Detail(ctx, pg, statusId)
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreateStatusReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.statusRepository.Create(tx, &entity.Status{
			Code: req.Code,
			Name: req.Name,
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, statusId uint, req *dto.UpdateStatusReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.statusRepository.Update(tx, &entity.Status{
			Model: gorm.Model{
				ID: statusId,
			},
			Code: req.Code,
			Name: req.Name,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, statusId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.statusRepository.Delete(tx, statusId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, statusRepository StatusRepository) *Service {
	return &Service{
		pg:               pg,
		auditRepository:  auditRepository,
		statusRepository: statusRepository,
	}
}
