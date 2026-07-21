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

type LanguageRepository interface {
	List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.LanguageListResponse, error)
	Detail(ctx context.Context, pg *sql.DB, languageId uint) (*readmodel.LanguageDetail, error)
	Create(tx *gorm.DB, language *entity.Language) error
	Update(tx *gorm.DB, language *entity.Language) error
	Delete(tx *gorm.DB, languageId uint) error
}

type Service struct {
	pg                 *gorm.DB
	auditRepository    AuditRepository
	languageRepository LanguageRepository
}

func (s *Service) List(ctx context.Context, page uint) (*readmodel.LanguageListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.languageRepository.List(ctx, pg, page)
}

func (s *Service) Detail(ctx context.Context, languageId uint) (*readmodel.LanguageDetail, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	return s.languageRepository.Detail(ctx, pg, languageId)
}

func (s *Service) Create(ctx context.Context, userId uint, req *dto.CreateLanguageReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.languageRepository.Create(tx, &entity.Language{
			Code: req.Code,
			Name: req.Name,
		})
	})
}

func (s *Service) Update(ctx context.Context, userId uint, languageId uint, req *dto.UpdateLanguageReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.languageRepository.Update(tx, &entity.Language{
			Model: gorm.Model{
				ID: languageId,
			},
			Code: req.Code,
			Name: req.Name,
		})
	})
}

func (s *Service) Delete(ctx context.Context, userId uint, languageId uint) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		return s.languageRepository.Delete(tx, languageId)
	})
}

func New(pg *gorm.DB, auditRepository AuditRepository, languageRepository LanguageRepository) *Service {
	return &Service{
		pg:                 pg,
		auditRepository:    auditRepository,
		languageRepository: languageRepository,
	}
}
