package service

import (
	"context"
	"database/sql"
	"fmt"

	"enterprise_resource_planning/services/material_management/dto"
	"enterprise_resource_planning/services/material_management/entity"
	"enterprise_resource_planning/services/material_management/readmodel"

	"gorm.io/gorm"
)

type AuditRepository interface {
	SetUserId(tx *gorm.DB, userId uint) error
}

type WarehouseRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.WarehouseListResponse, error)
	Create(tx *gorm.DB, warehouse *entity.Warehouse) error
	Update(tx *gorm.DB, warehouse *entity.Warehouse) error
	Delete(tx *gorm.DB, warehouseId uint) error
}

type Service struct {
	pg                  *gorm.DB
	auditRepository     AuditRepository
	warehouseRepository WarehouseRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.WarehouseListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	data, err := s.warehouseRepository.List(ctx, pg, page)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreateWarehouseReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.warehouseRepository.Create(tx, &entity.Warehouse{
			Name:        req.Name,
			Description: req.Description,
			Address:     req.Address,
			Location:    fmt.Sprintf("SRID=4326;POINT(%f %f)", req.Location.Longitude, req.Location.Latitude),
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, warehouseId uint, req *dto.UpdateWarehouseReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.warehouseRepository.Update(tx, &entity.Warehouse{
			Model: gorm.Model{
				ID: warehouseId,
			},
			Name:        req.Name,
			Description: req.Description,
			Address:     req.Address,
			Location:    fmt.Sprintf("SRID=4326;POINT(%f %f)", req.Location.Longitude, req.Location.Latitude),
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, warehouseId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.warehouseRepository.Delete(tx, warehouseId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, warehouseRepository WarehouseRepository) *Service {
	return &Service{
		pg:                  pg,
		auditRepository:     auditRepository,
		warehouseRepository: warehouseRepository,
	}
}
