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

type DepartmentRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.DepartmentListResponse, error)
	Detail(ctx context.Context, pg *sql.DB, departmentId uint) (*readmodel.DepartmentDetail, error)
	Create(tx *gorm.DB, mu *entity.Department) error
	Update(tx *gorm.DB, mu *entity.Department) error
	Delete(tx *gorm.DB, departmentId uint) error
}

type Service struct {
	pg                   *gorm.DB
	auditRepository      AuditRepository
	departmentRepository DepartmentRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.DepartmentListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.departmentRepository.List(ctx, pg, page)
}

func (s *Service) Detail(ctx context.Context, departmentId uint) (*readmodel.DepartmentDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.departmentRepository.Detail(ctx, pg, departmentId)
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreateDepartmentReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.departmentRepository.Create(tx, &entity.Department{
			Name:        req.Name,
			Description: req.Description,
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, departmentId uint, req *dto.UpdateDepartmentReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.departmentRepository.Update(tx, &entity.Department{
			Model: gorm.Model{
				ID: departmentId,
			},
			Name:        req.Name,
			Description: req.Description,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, departmentId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.departmentRepository.Delete(tx, departmentId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, departmentRepository DepartmentRepository) *Service {
	return &Service{
		pg:                   pg,
		auditRepository:      auditRepository,
		departmentRepository: departmentRepository,
	}
}
