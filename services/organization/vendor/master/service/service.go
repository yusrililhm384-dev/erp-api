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

type VendorMasterRepository interface {
	Detail(ctx context.Context, pg *sql.DB, vendorId uint) (*readmodel.VendorDetail, error)
	List(ctx context.Context, pg *sql.DB) (*readmodel.VendorListResponse, error)
	Create(tx *gorm.DB, vendor *entity.Vendor) error
	Update(tx *gorm.DB, vendor *entity.Vendor) error
	Delete(tx *gorm.DB, vendorId uint) error
}

type Service struct {
	pg                     *gorm.DB
	auditRepository        AuditRepository
	vendorMasterRepository VendorMasterRepository
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreateVendorReq) error {
	vendor := &entity.Vendor{
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		Location:    fmt.Sprintf("SRID=4326;POINT(%f %f)", req.Location.Longitude, req.Location.Latitude),
	}

	for _, c := range req.Contacts {
		vendor.Contacts = append(vendor.Contacts, &entity.VendorContact{
			Name:  c.Name,
			Email: c.Email,
			Phone: c.Phone,
		})
	}

	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.vendorMasterRepository.Create(tx, vendor)
	})
}

func (s *Service) Update(ctx context.Context, userId uint, vendorId uint, req *dto.UpdateVendorReq) error {
	vendor := &entity.Vendor{
		Model: gorm.Model{
			ID: vendorId,
		},
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		Location:    fmt.Sprintf("SRID=4326;POINT(%f %f)", req.Location.Longitude, req.Location.Latitude),
	}

	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.vendorMasterRepository.Update(tx, vendor)
	})
}

func (s *Service) Delete(ctx context.Context, userId, vendorId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.vendorMasterRepository.Delete(tx, vendorId)
	})
}

func (s *Service) Detail(ctx context.Context, vendorId uint) (*readmodel.VendorDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.vendorMasterRepository.Detail(ctx, pg, vendorId)
}

func (s *Service) List(ctx context.Context) (*readmodel.VendorListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.vendorMasterRepository.List(ctx, pg)
}

func New(pg *gorm.DB, auditRepository AuditRepository, vendorMasterRepository VendorMasterRepository) *Service {
	return &Service{
		pg:                     pg,
		auditRepository:        auditRepository,
		vendorMasterRepository: vendorMasterRepository,
	}
}
