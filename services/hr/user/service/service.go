package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"enterprise_resource_planning/services/hr/dto"
	"enterprise_resource_planning/services/hr/entity"

	errs "enterprise_resource_planning/internal/shared/errors"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuditRepository interface {
	SetUserId(tx *gorm.DB, userId uint) error
}

type UserRepository interface {
	Update(tx *gorm.DB, user *entity.User) error
	FindByUsername(ctx context.Context, pg *sql.DB, username string) (*entity.User, error)
}

type CacheRepository interface {
	SetJSON(ctx context.Context, key string, val any, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	IncrExpired(ctx context.Context, key string, ttl time.Duration) (int64, error)
	Set(ctx context.Context, key string, val string, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	SetNX(ctx context.Context, key string, val any, ttl time.Duration) error
}

type Service struct {
	pg              *gorm.DB
	auditRepository AuditRepository
	userRepository  UserRepository
	cacheRepository CacheRepository
}

func (s *Service) Logout(ctx context.Context, sessionId string) error {
	if err := s.cacheRepository.Delete(ctx, sessionId); err != nil {
		if errors.Is(err, redis.Nil) {
			return errs.ErrInvalidSession
		}

		return err
	}

	return nil
}

func (s *Service) ChangePassword(ctx context.Context, userId uint, req *dto.ChangePasswordReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		u := &entity.User{}

		if err := tx.Model(&entity.User{}).Where("employee_id = ?", userId).First(u).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errs.ErrUserNotFound
			}

			return err
		}

		if err := u.CompareHashAndPassword(req.OldPassword); err != nil {
			return err
		}

		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		user := &entity.User{
			Password:   req.NewPassword,
			EmployeeID: userId,
		}

		if err := user.GenerateHashPassowrd(); err != nil {
			return err
		}

		return s.userRepository.Update(tx, user)
	})
}

func (s *Service) ChangeUsername(ctx context.Context, userId uint, req *dto.ChangeUsernameReq) error {
	return s.pg.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.auditRepository.SetUserId(tx, userId); err != nil {
			return err
		}

		user := &entity.User{
			Username:   req.Username,
			EmployeeID: userId,
		}

		return s.userRepository.Update(tx, user)
	})
}

func (s *Service) Login(ctx context.Context, req *dto.UserLoginReq) (*dto.UserLoginResponse, error) {
	lockKey, attemptKey := fmt.Sprintf("locked:%s", req.Username), fmt.Sprintf("attempt:%s", req.Username)

	res, err := s.cacheRepository.Get(ctx, lockKey)

	if res != "" && err == nil {
		return nil, errs.ErrUserLocked
	}

	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	data, err := s.userRepository.FindByUsername(ctx, pg, req.Username)

	if err != nil {
		return nil, err
	}

	if err := data.CompareHashAndPassword(req.Password); err != nil {
		count, err := s.cacheRepository.IncrExpired(ctx, attemptKey, 15*time.Minute)

		if err != nil {
			return nil, err
		}

		if count >= 3 {
			if err := s.cacheRepository.SetNX(ctx, lockKey, "locked", 15*time.Minute); err != nil {
				return nil, err
			}

			_ = s.cacheRepository.Delete(ctx, attemptKey)

			return nil, errs.ErrTooManyRequest
		}

		return nil, errs.ErrInvalidPassword
	}

	sessionId := fmt.Sprintf("session_id:%s", uuid.NewString())

	if err := s.cacheRepository.SetJSON(ctx, sessionId, &dto.UserSession{
		SessionId: sessionId,
		UserId:    data.ID,
		Username:  data.Username,
		Role:      data.Role,
	}, 8*time.Hour); err != nil {
		return nil, err
	}

	_ = s.cacheRepository.Delete(ctx, lockKey, attemptKey)

	return &dto.UserLoginResponse{
		SessionId: sessionId,
	}, nil
}

func New(pg *sql.DB, auditRepository AuditRepository, userRepository UserRepository, cacheRepository CacheRepository) *Service {
	return &Service{
		pg:              &gorm.DB{},
		auditRepository: auditRepository,
		userRepository:  userRepository,
		cacheRepository: cacheRepository,
	}
}
