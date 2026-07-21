package service

import (
	"context"
	"time"

	"enterprise_resource_planning/services/hr/dto"
	"enterprise_resource_planning/services/hr/entity"

	"gorm.io/gorm"
)

type AuditRepository interface {
	SetUserId(tx *gorm.DB, userId uint) error
}

type MutationRepository interface {
	Create(tx *gorm.DB, mutaion *entity.EmployeeMutation) error
	Update(tx *gorm.DB, mutaion *entity.EmployeeMutation) error
}

type Service struct {
	pg                 *gorm.DB
	auditRepository    AuditRepository
	mutationRepository MutationRepository
}

func (s *Service) Mutation(ctx context.Context, userId uint, employeeId uint, req *dto.CreateMutationReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		beforeStartDate := req.StartDate.AddDate(0, 0, -1)

		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		if err := s.mutationRepository.Update(tx, &entity.EmployeeMutation{
			Model: gorm.Model{
				DeletedAt: gorm.DeletedAt{
					Time:  time.Now(),
					Valid: true,
				},
				UpdatedAt: time.Now(),
			},
			EmployeeID: employeeId,
			EndDate:    &beforeStartDate,
		}); err != nil {
			return err
		}

		return s.mutationRepository.Create(tx, &entity.EmployeeMutation{
			EmployeeID:   employeeId,
			BranchID:     req.BranchID,
			PositionID:   req.PositionID,
			TypeID:       req.TypeID,
			DepartmentID: req.DepartmentID,
			StartDate:    req.StartDate,
			EndDate:      req.EndDate,
		})
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, mutationRepository MutationRepository) *Service {
	return &Service{
		pg:                 pg,
		auditRepository:    auditRepository,
		mutationRepository: mutationRepository,
	}
}
