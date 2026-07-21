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

type TypeRepository interface {
	Create(tx *gorm.DB, t *entity.EmployeeType) error
	Update(tx *gorm.DB, t *entity.EmployeeType) error
	Delete(tx *gorm.DB, typeId uint) error
	Detail(ctx context.Context, pg *sql.DB, typeId uint) (*entity.EmployeeType, error)
	List(ctx context.Context, pg *sql.DB) ([]*readmodel.Type, error)
}

type Service struct {
	pg              *gorm.DB
	auditRepository AuditRepository
	typeRepository  TypeRepository
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreateTypeReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.typeRepository.Create(tx, &entity.EmployeeType{
			Name:        req.Name,
			Description: req.Description,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId, typeId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.typeRepository.Delete(tx, typeId)
	})
}

func (s *Service) Update(ctx context.Context, userId uint, typeId uint, req *dto.UpdateTypeReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.typeRepository.Update(tx, &entity.EmployeeType{
			Model: gorm.Model{
				ID: typeId,
			},
			Name:        req.Name,
			Description: req.Description,
		})
	})
}

func (s *Service) Detail(ctx context.Context, typeId uint) (*dto.DetailTypeRes, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	data, err := s.typeRepository.Detail(ctx, pg, typeId)

	if err != nil {
		return nil, err
	}

	return &dto.DetailTypeRes{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}, nil
}

func (s *Service) List(ctx context.Context) ([]*readmodel.Type, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	data, err := s.typeRepository.List(ctx, pg)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func New(pg *gorm.DB, auditRepository AuditRepository, typeRepository TypeRepository) *Service {
	return &Service{
		pg:              pg,
		auditRepository: auditRepository,
		typeRepository:  typeRepository,
	}
}
