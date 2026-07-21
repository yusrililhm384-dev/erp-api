package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"enterprise_resource_planning/internal/shared/response"
	"enterprise_resource_planning/services/hr/dto"

	errs "enterprise_resource_planning/internal/shared/errors"
	auth "enterprise_resource_planning/services/auth/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type Service interface {
	Mutation(ctx context.Context, userId uint, employeeId uint, req *dto.CreateMutationReq) error
}

type Handler struct {
	service  Service
	validate *validator.Validate
}

func (h *Handler) Mutation(c *echo.Context) error {
	c.Logger().Info("Employee Mutation")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	employeeId, err := strconv.ParseUint(c.Param("employeeId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid employee id")
	}

	req := &dto.CreateMutationReq{}

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid body request!")
	}

	if err := h.validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing require fields!")
	}

	if err := h.service.Mutation(c.Request().Context(), sessionData.UserId, uint(employeeId), req); err != nil {
		c.Logger().Warn(err.Error())

		switch {
		case errors.Is(err, errs.ErrForeignKeyError):
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid reference id!")
		case errors.Is(err, errs.ErrEmployeeMaserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Employee master not found!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	return c.JSON(http.StatusCreated, &response.Response{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Employee successfully mutation!",
	})
}

func New(service Service, validate *validator.Validate) *Handler {
	return &Handler{
		service:  service,
		validate: validate,
	}
}
