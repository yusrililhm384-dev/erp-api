package service

import (
	"context"
	"database/sql"

	"enterprise_resource_planning/services/hr/dto"
	"enterprise_resource_planning/services/hr/entity"
	"enterprise_resource_planning/services/hr/readmodel"

	"gorm.io/gorm"
)

type AuditRepository interface {
	SetUserId(tx *gorm.DB, userId uint) error
}

type MasterRepository interface {
	Create(tx *gorm.DB, employee *entity.Employee) error
	Update(tx *gorm.DB, employee *entity.Employee) error
	Delete(tx *gorm.DB, employeeId uint) error
	Detail(ctx context.Context, pg *sql.DB, employeeId uint) (*readmodel.EmployeeDetail, error)
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.EmployeeListResponse, error)
}

type LegalRepository interface {
	Create(tx *gorm.DB, legal *entity.EmployeeLegal) error
}

type MutationRepository interface {
	Create(tx *gorm.DB, mutation *entity.EmployeeMutation) error
}

type UserRepository interface {
	Create(tx *gorm.DB, user *entity.User) error
}

type Service struct {
	pg                 *gorm.DB
	auditRepository    AuditRepository
	masterRepository   MasterRepository
	legalRepository    LegalRepository
	mutationRepository MutationRepository
	userRepository     UserRepository
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreateEmployeeReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		employee := &entity.Employee{
			Name:      req.Name,
			Email:     req.Email,
			Phone:     req.Phone,
			Address:   req.Address,
			BirthDate: req.BirthDate,
			JoinDate:  req.JoinDate,
		}

		if err := s.masterRepository.Create(tx, employee); err != nil {
			return err
		}

		if err := s.legalRepository.Create(tx, &entity.EmployeeLegal{
			EmployeeID:            employee.ID,
			IdentityNumber:        employee.Legal.IdentityNumber,
			TaxNumber:             req.Legal.TaxNumber,
			HealthInsuranceNumber: employee.Legal.HealthInsuranceNumber,
			LaborInsuranceNumber:  employee.Legal.LaborInsuranceNumber,
		}); err != nil {
			return err
		}

		if err := s.mutationRepository.Create(tx, &entity.EmployeeMutation{
			EmployeeID: employee.ID,
			BranchID:   req.Mutation.BranchID,
			PositionID: req.Mutation.PositionID,
			StartDate:  req.Mutation.StartDate,
			TypeID:     req.Mutation.TypeID,
			EndDate:    req.Mutation.EndDate,
		}); err != nil {
			return err
		}

		user := &entity.User{
			Username:   req.TempUsername,
			Password:   req.TempPassword,
			EmployeeID: employee.ID,
		}

		if err := user.GenerateHashPassowrd(); err != nil {
			return err
		}

		return s.userRepository.Create(tx, user)
	})
}

func (s *Service) Update(ctx context.Context, employeeId uint, req *dto.UpdateEmployeeReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, employeeId); err != nil {
			return err
		}

		return s.masterRepository.Update(tx, &entity.Employee{
			Model: gorm.Model{
				ID: employeeId,
			},
			Name:      req.Name,
			Email:     req.Email,
			Phone:     req.Phone,
			Address:   req.Address,
			BirthDate: req.BirthDate,
			JoinDate:  req.JoinDate,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId, employeeId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.masterRepository.Delete(tx, employeeId)
	})
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.EmployeeListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	data, err := s.masterRepository.List(ctx, pg, page)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Service) Detail(ctx context.Context, employeeId uint) (*readmodel.EmployeeDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	data, err := s.masterRepository.Detail(ctx, pg, employeeId)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func New(pg *gorm.DB, auditRepository AuditRepository, masterRepository MasterRepository, legalRepository LegalRepository, mutationRepository MutationRepository, userRepository UserRepository) *Service {
	return &Service{
		pg:                 pg,
		auditRepository:    auditRepository,
		masterRepository:   masterRepository,
		legalRepository:    legalRepository,
		mutationRepository: mutationRepository,
		userRepository:     userRepository,
	}
}
