package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"enterprise_resource_planning/services/material_management/dto"
	"enterprise_resource_planning/services/material_management/entity"
	"enterprise_resource_planning/services/material_management/readmodel"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

type CacheRepository interface {
	Delete(ctx context.Context, keys ...string) error
	GetJSON(ctx context.Context, key string, dest any) error
	SetJSON(ctx context.Context, key string, val any, ttl time.Duration) error
	Incr(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val string, ttl time.Duration) error
}

type AuditRepository interface {
	SetUserId(tx *gorm.DB, userId uint) error
}

type MaterialMasterRepository interface {
	Create(tx *gorm.DB, mm *entity.MaterialMaster) error
	Update(tx *gorm.DB, mm *entity.MaterialMaster) error
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.MaterialMasterListResponse, error)
}

type Service struct {
	pg                       *gorm.DB
	auditRepository          AuditRepository
	materialMasterRepository MaterialMasterRepository
	cacheRepository          CacheRepository
	group                    *singleflight.Group
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.MaterialMasterListResponse, error) {
	const versionKey = "material_versions"

	version, err := s.cacheRepository.Get(ctx, versionKey)
	if err != nil {
		version = "1"
		_ = s.cacheRepository.Set(ctx, versionKey, version, 0)
	}

	key := fmt.Sprintf("material:v%s:page:%d", version, page)

	pg, err := s.pg.DB()
	if err != nil {
		return nil, err
	}

	val, err, _ := s.group.Do(key, func() (any, error) {
		cached := new(readmodel.MaterialMasterListResponse)

		if err := s.cacheRepository.GetJSON(ctx, key, cached); err == nil {
			return cached, nil
		}

		result, err := s.materialMasterRepository.List(ctx, pg, page)
		if err != nil {
			return nil, err
		}

		_ = s.cacheRepository.SetJSON(ctx, key, result, 15*time.Minute)

		return result, nil
	})

	if err != nil {
		return nil, err
	}

	return val.(*readmodel.MaterialMasterListResponse), nil
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreateMaterialMasterReq) error {
	if err := s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		mm := &entity.MaterialMaster{
			Code:                   req.Code,
			Name:                   req.Name,
			Description:            req.Description,
			Barcode:                req.Barcode,
			MaterialTypeID:         req.MaterialTypeId,
			MaterialGroupID:        req.MaterialGroupId,
			ManufacturerID:         req.ManufacturerId,
			ManufacturerPartNumber: req.ManufacturerPartNumber,
			BrandID:                req.BrandId,
			BaseUnitID:             req.BaseUnitId,
			CountryID:              req.CountryId,
			MaterialStatusID:       req.MaterialStatusId,
		}

		if req.MaterialPhysical != nil {
			mm.MaterialPhysical = &entity.MaterialPhysical{
				WeightUnitID:    req.MaterialPhysical.WeightUnitId,
				GrossWeight:     req.MaterialPhysical.GrossWeight,
				NetWeight:       req.MaterialPhysical.NetWeight,
				VolumeUnitID:    req.MaterialPhysical.VolumeUnitId,
				Volume:          req.MaterialPhysical.Volume,
				DimensionUnitID: req.MaterialPhysical.DimensionUnitId,
				Length:          req.MaterialPhysical.Length,
				Width:           req.MaterialPhysical.Width,
				Height:          req.MaterialPhysical.Height,
				IsSerialManaged: req.MaterialPhysical.IsSerialManaged,
				IsBatchManaged:  req.MaterialPhysical.IsBatchManaged,
			}
		}

		return s.materialMasterRepository.Create(tx, mm)
	}); err != nil {
		return err
	}

	_ = s.cacheRepository.Incr(ctx, "material_versions")

	return nil
}

func (s *Service) Update(ctx context.Context, userId uint, materialMasterId uint, req *dto.UpdateMaterialMasterReq) error {
	if err := s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		mm := &entity.MaterialMaster{
			Model: gorm.Model{
				ID: materialMasterId,
			},
			Code:                   req.Code,
			Name:                   req.Name,
			Description:            req.Description,
			Barcode:                req.Barcode,
			MaterialTypeID:         req.MaterialTypeId,
			MaterialGroupID:        req.MaterialGroupId,
			ManufacturerID:         req.ManufacturerId,
			ManufacturerPartNumber: req.ManufacturerPartNumber,
			BrandID:                req.BrandId,
			BaseUnitID:             req.BaseUnitId,
			CountryID:              req.CountryId,
			MaterialStatusID:       req.MaterialStatusId,
		}

		return s.materialMasterRepository.Update(tx, mm)
	}); err != nil {
		return err
	}

	_ = s.cacheRepository.Delete(ctx, fmt.Sprintf("material:%d", materialMasterId))
	_ = s.cacheRepository.Incr(ctx, "material_versions")

	return nil
}

func New(pg *gorm.DB, auditRepository AuditRepository, materialMasterRepository MaterialMasterRepository, cacheRepository CacheRepository) *Service {
	return &Service{
		pg:                       pg,
		auditRepository:          auditRepository,
		materialMasterRepository: materialMasterRepository,
		cacheRepository:          cacheRepository,
	}
}
