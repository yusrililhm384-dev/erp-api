package service

import (
	"context"

	"enterprise_resource_planning/services/hr/dto"
	"enterprise_resource_planning/services/hr/entity"

	"gorm.io/gorm"
)

type AuditRepository interface {
	SetUserId(tx *gorm.DB, userId uint) error
}

type LegalRepository interface {
	Update(tx *gorm.DB, legal *entity.EmployeeLegal) error
}

type Service struct {
	pg              *gorm.DB
	auditRepository AuditRepository
	legalRepository LegalRepository
}

func (s *Service) Update(ctx context.Context, userId uint, req *dto.UpdateLegalReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.legalRepository.Update(tx, &entity.EmployeeLegal{
			EmployeeID:            1,
			IdentityNumber:        req.IdentityNumber,
			TaxNumber:             req.TaxNumber,
			HealthInsuranceNumber: req.HealthInsuranceNumber,
			LaborInsuranceNumber:  req.LaborInsuranceNumber,
		})
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, legalRepository LegalRepository) *Service {
	return &Service{
		pg:              pg,
		auditRepository: auditRepository,
		legalRepository: legalRepository,
	}
}
