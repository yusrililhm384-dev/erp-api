package service

import (
	"context"
	"fmt"
	"net/http"

	"enterprise_resource_planning/internal/shared/constant"
	"enterprise_resource_planning/internal/shared/role"
	"enterprise_resource_planning/services/hr/dto"

	"github.com/labstack/echo/v5"
)

type CacheRepository interface {
	GetJSON(ctx context.Context, key string, dest any) error
}

type Service struct {
	cacheRepository CacheRepository
}

func GetSessionFromContext(c *echo.Context) (*dto.UserSession, error) {
	sessionData, ok := c.Request().Context().Value(constant.SessionKey).(*dto.UserSession)

	if !ok {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	return sessionData, nil
}

func (s *Service) Authentication() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			cookie, err := c.Cookie("session_id")

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "You must login first")
			}

			sessionId := cookie.Value

			if sessionId == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "You must login first!")
			}

			userSession := &dto.UserSession{}

			if err := s.cacheRepository.GetJSON(c.Request().Context(), sessionId, userSession); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Session expired!")
			}

			ctx := context.WithValue(c.Request().Context(), constant.SessionKey, userSession)

			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func (s *Service) Authorization(resource, action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			data, err := GetSessionFromContext(c)

			if err != nil {
				return err
			}

			if !s.Can(data.Role, resource, action) {
				return echo.NewHTTPError(http.StatusForbidden, "You're not authorized!")
			}

			return next(c)
		}
	}
}

func (s *Service) Can(userRole, resource, action string) bool {
	key := fmt.Sprintf("%s.%s", resource, action)

	permission, ok := role.Permission[userRole]

	if !ok {
		return false
	}

	return permission[key]
}

func New(cacheRepository CacheRepository) *Service {
	return &Service{
		cacheRepository: cacheRepository,
	}
}
