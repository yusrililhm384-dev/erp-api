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

type ContactRepository interface {
	Create(tx *gorm.DB, contact *entity.VendorContact) error
	Update(tx *gorm.DB, contact *entity.VendorContact) error
	Delete(tx *gorm.DB, contactId, vendorId uint) error
	Detail(ctx context.Context, pg *sql.DB, contactId, vendorId uint) (*readmodel.VendorContactDetail, error)
}

type Service struct {
	pg                *gorm.DB
	auditRepository   AuditRepository
	contactRepository ContactRepository
}

func (s *Service) Create(ctx context.Context, userId uint, vendorId uint, req *dto.ContactVendorReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.contactRepository.Create(tx, &entity.VendorContact{
			Name:     req.Name,
			Email:    req.Email,
			Phone:    req.Phone,
			VendorID: vendorId,
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, vendorId uint, contactId uint, req *dto.UpdateContactVendorReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.contactRepository.Update(tx, &entity.VendorContact{
			Model: gorm.Model{
				ID: contactId,
			},
			Name:     req.Name,
			Email:    req.Email,
			Phone:    req.Phone,
			VendorID: vendorId,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, vendorId uint, contactId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.contactRepository.Delete(tx, contactId, vendorId)
	})
}

func (s *Service) Detail(ctx context.Context, contactId, vendorId uint) (*readmodel.VendorContactDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.contactRepository.Detail(ctx, pg, contactId, vendorId)
}

func New(pg *gorm.DB, auditRepository AuditRepository, contactRepository ContactRepository) *Service {
	return &Service{
		pg:                pg,
		auditRepository:   auditRepository,
		contactRepository: contactRepository,
	}
}
