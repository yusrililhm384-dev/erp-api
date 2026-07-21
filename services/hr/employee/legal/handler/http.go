package handler

import (
	"context"
	"errors"
	"net/http"

	"enterprise_resource_planning/internal/shared/response"
	"enterprise_resource_planning/services/hr/dto"

	errs "enterprise_resource_planning/internal/shared/errors"
	auth "enterprise_resource_planning/services/auth/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type Service interface {
	Update(ctx context.Context, userId uint, req *dto.UpdateLegalReq) error
}

type Handler struct {
	service  Service
	validate *validator.Validate
}

func (h *Handler) Update(c *echo.Context) error {
	c.Logger().Info("Update Legal")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	req := &dto.UpdateLegalReq{}

	if err := c.Bind(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid body request")
	}

	if err := h.validate.Struct(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid body request")
	}

	if err := h.service.Update(c.Request().Context(), sessionData.UserId, req); err != nil {
		c.Logger().Warn(err.Error())

		switch {
		case errors.Is(err, errs.ErrDuplicateLegalError):
			return echo.NewHTTPError(http.StatusConflict, "Duplicate employee legal!")
		case errors.Is(err, errs.ErrEmployeeMaserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Master not found!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Legal successfully updated!",
	})
}

func New(service Service, validate *validator.Validate) *Handler {
	return &Handler{
		service:  service,
		validate: validate,
	}
}
