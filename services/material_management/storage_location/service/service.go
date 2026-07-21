package service

import (
	"context"
	"database/sql"

	"enterprise_resource_planning/services/material_management/dto"
	"enterprise_resource_planning/services/material_management/entity"
	"enterprise_resource_planning/services/material_management/readmodel"

	"gorm.io/gorm"
)

type StorageLocationRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.StorageLocationListResponse, error)
	Detail(ctx context.Context, pg *sql.DB, storageLocationId uint) (*readmodel.StorageLocationDetail, error)
	Create(tx *gorm.DB, sl *entity.StorageLocation) error
	Update(tx *gorm.DB, sl *entity.StorageLocation) error
	Delete(tx *gorm.DB, storageLocationId uint) error
}

type AuditRepository interface {
	SetUserId(tx *gorm.DB, userId uint) error
}

type Service struct {
	pg                        *gorm.DB
	auditRepository           AuditRepository
	storageLocationRepository StorageLocationRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.StorageLocationListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.storageLocationRepository.List(ctx, pg, page)
}

func (s *Service) Detail(ctx context.Context, storageLocationId uint) (*readmodel.StorageLocationDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.storageLocationRepository.Detail(ctx, pg, storageLocationId)
}

func (s *Service) Create(ctx context.Context, userId uint, warehouseId uint, req *dto.CreateStorageLocationReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.storageLocationRepository.Create(tx, &entity.StorageLocation{
			Code:        req.Code,
			Name:        req.Name,
			Description: req.Description,
			WarehouseID: warehouseId,
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, storageLocationId uint, req *dto.UpdateStorageLocationReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.storageLocationRepository.Update(tx, &entity.StorageLocation{
			Model: gorm.Model{
				ID: storageLocationId,
			},
			Name:        req.Name,
			Description: req.Description,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, storageLocationId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.storageLocationRepository.Delete(tx, storageLocationId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, storageLocationRepository StorageLocationRepository) *Service {
	return &Service{
		pg:                        pg,
		auditRepository:           auditRepository,
		storageLocationRepository: storageLocationRepository,
	}
}
