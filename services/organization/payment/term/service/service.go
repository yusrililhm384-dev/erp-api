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

type PaymentTermRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.PaymentTermListResponse, error)
	Detail(ctx context.Context, pg *sql.DB, paymentTermId uint) (*readmodel.PaymentTermDetail, error)
	Create(tx *gorm.DB, pm *entity.PaymentTerm) error
	Update(tx *gorm.DB, pm *entity.PaymentTerm) error
	Delete(tx *gorm.DB, paymentTermId uint) error
}

type Service struct {
	pg                    *gorm.DB
	auditRepository       AuditRepository
	paymentTermRepository PaymentTermRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.PaymentTermListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.paymentTermRepository.List(ctx, pg, page)
}

func (s *Service) Detail(ctx context.Context, paymentTermId uint) (*readmodel.PaymentTermDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.paymentTermRepository.Detail(ctx, pg, paymentTermId)
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreatePaymentTermReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.paymentTermRepository.Create(tx, &entity.PaymentTerm{
			Code:        req.Code,
			Name:        req.Name,
			Description: req.Description,
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, paymentTermId uint, req *dto.UpdatePaymentTermReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.paymentTermRepository.Update(tx, &entity.PaymentTerm{
			Model: gorm.Model{
				ID: paymentTermId,
			},
			Code:        req.Code,
			Name:        req.Name,
			Description: req.Description,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, paymentTermId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.paymentTermRepository.Delete(tx, paymentTermId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, paymentTermRepository PaymentTermRepository) *Service {
	return &Service{
		pg:                    pg,
		auditRepository:       auditRepository,
		paymentTermRepository: paymentTermRepository,
	}
}
