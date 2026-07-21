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

type ManufacturerRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.ManufacturerListResponse, error)
	Detail(ctx context.Context, pg *sql.DB, manufacturerId uint) (*readmodel.ManufacturerDetail, error)
	Create(tx *gorm.DB, mf *entity.Manufacturer) error
	Update(tx *gorm.DB, mf *entity.Manufacturer) error
	Delete(tx *gorm.DB, manufacturerId uint) error
}

type Service struct {
	pg                     *gorm.DB
	auditRepository        AuditRepository
	manufacturerRepository ManufacturerRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.ManufacturerListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.manufacturerRepository.List(ctx, pg, page)
}

func (s *Service) Detail(ctx context.Context, manufacturerId uint) (*readmodel.ManufacturerDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.manufacturerRepository.Detail(ctx, pg, manufacturerId)
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreateManufacturerReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.manufacturerRepository.Create(tx, &entity.Manufacturer{
			Code: req.Code,
			Name: req.Name,
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, manufacturerId uint, req *dto.UpdateManufacturerReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.manufacturerRepository.Update(tx, &entity.Manufacturer{
			Model: gorm.Model{
				ID: manufacturerId,
			},
			Code: req.Code,
			Name: req.Name,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, manufacturerId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.manufacturerRepository.Delete(tx, manufacturerId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, manufacturerRepository ManufacturerRepository) *Service {
	return &Service{
		pg:                     pg,
		auditRepository:        auditRepository,
		manufacturerRepository: manufacturerRepository,
	}
}
