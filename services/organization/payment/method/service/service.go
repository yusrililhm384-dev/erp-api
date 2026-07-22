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

type PaymentMethodRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.PaymentMethodListResponse, error)
	Detail(ctx context.Context, pg *sql.DB, paymentMethodId uint) (*readmodel.PaymentMethodDetail, error)
	Create(tx *gorm.DB, pm *entity.PaymentMethod) error
	Update(tx *gorm.DB, pm *entity.PaymentMethod) error
	Delete(tx *gorm.DB, paymentMethodId uint) error
}

type Service struct {
	pg                      *gorm.DB
	auditRepository         AuditRepository
	paymentMethodRepository PaymentMethodRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.PaymentMethodListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.paymentMethodRepository.List(ctx, pg, page)
}

func (s *Service) Detail(ctx context.Context, paymentMethodId uint) (*readmodel.PaymentMethodDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.paymentMethodRepository.Detail(ctx, pg, paymentMethodId)
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreatePaymentMethodReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.paymentMethodRepository.Create(tx, &entity.PaymentMethod{
			Code:        req.Code,
			Name:        req.Name,
			Description: req.Description,
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, paymentMethodId uint, req *dto.UpdatePaymentMethodReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.paymentMethodRepository.Update(tx, &entity.PaymentMethod{
			Model: gorm.Model{
				ID: paymentMethodId,
			},
			Code:        req.Code,
			Name:        req.Name,
			Description: req.Description,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, paymentMethodId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.paymentMethodRepository.Delete(tx, paymentMethodId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, paymentMethodRepository PaymentMethodRepository) *Service {
	return &Service{
		pg:                      pg,
		auditRepository:         auditRepository,
		paymentMethodRepository: paymentMethodRepository,
	}
}
