package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"enterprise_resource_planning/internal/shared/response"
	"enterprise_resource_planning/services/hr/dto"
	"enterprise_resource_planning/services/hr/readmodel"

	errs "enterprise_resource_planning/internal/shared/errors"
	auth "enterprise_resource_planning/services/auth/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type Service interface {
	Create(ctx context.Context, userId uint, req *dto.CreateEmployeeReq) error
	Update(ctx context.Context, employeeId uint, req *dto.UpdateEmployeeReq) error
	Delete(ctx context.Context, userId, employeeId uint) error
	Detail(ctx context.Context, employeeId uint) (*readmodel.EmployeeDetail, error)
	List(ctx context.Context, page uint) (*readmodel.EmployeeListResponse, error)
}

type Handler struct {
	service  Service
	validate *validator.Validate
}

func (h *Handler) Create(c *echo.Context) error {
	c.Logger().Info("Create Employee")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	req := &dto.CreateEmployeeReq{}

	if err := c.Bind(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid body request!")
	}

	if err := h.validate.Struct(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadGateway, "Missing required request!")
	}

	if err := h.service.Create(c.Request().Context(), sessionData.UserId, req); err != nil {
		c.Logger().Warn(err.Error())

		switch {
		case errors.Is(err, errs.ErrDuplicateContactError):
			return echo.NewHTTPError(http.StatusConflict, "Contact already exists!")
		case errors.Is(err, errs.ErrForeignKeyError):
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid reference id!")
		case errors.Is(err, errs.ErrDuplicateLegalError):
			return echo.NewHTTPError(http.StatusConflict, "Legal already exists!")
		case errors.Is(err, errs.ErrDuplicateUsernameError):
			return echo.NewHTTPError(http.StatusConflict, "Username already exists!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	return c.JSON(http.StatusCreated, &response.Response{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Employee successfully created!",
	})
}

func (h *Handler) Update(c *echo.Context) error {
	c.Logger().Info("Update Employee")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	req := &dto.UpdateEmployeeReq{}

	if err := c.Bind(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid body request!")
	}

	if err := h.validate.Struct(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadGateway, "Missing required request!")
	}

	if err := h.service.Update(c.Request().Context(), sessionData.UserId, req); err != nil {
		c.Logger().Warn(err.Error())

		switch {
		case errors.Is(err, errs.ErrDuplicateContactError):
			return echo.NewHTTPError(http.StatusConflict, "Contact already exists!")
		case errors.Is(err, errs.ErrForeignKeyError):
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid reference id!")
		case errors.Is(err, errs.ErrEmployeeNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Employee not found!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Employee successfully updated!",
	})
}

func (h *Handler) Delete(c *echo.Context) error {
	c.Logger().Info("Delete Employee")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	employeeId, err := strconv.ParseUint(c.Param("employeeId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid employee id!")
	}

	if err := h.service.Delete(c.Request().Context(), sessionData.UserId, uint(employeeId)); err != nil {
		c.Logger().Warn(err.Error())

		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "User not found!")
		case errors.Is(err, errs.ErrEmployeeNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Employee not found!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Employee successfully deleted!",
	})
}

func (h *Handler) Profile(c *echo.Context) error {
	c.Logger().Info("Employee Profile")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	data, err := h.service.Detail(c.Request().Context(), sessionData.UserId)

	if err != nil {
		c.Logger().Warn(err.Error())

		if errors.Is(err, errs.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found!")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Profile successfully fetched!",
		Data:    data,
	})
}

func (h *Handler) Detail(c *echo.Context) error {
	c.Logger().Info("Employee Detail")

	employeeId, err := strconv.ParseUint(c.Param("employeeId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid employee id!")
	}

	data, err := h.service.Detail(c.Request().Context(), uint(employeeId))

	if err != nil {
		c.Logger().Warn(err.Error())

		if errors.Is(err, errs.ErrEmployeeNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Employee not found!")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Employee successfully fetched!",
		Data:    data,
	})
}

func (h *Handler) List(c *echo.Context) error {
	c.Logger().Info("Employee List")

	page, err := strconv.ParseUint(c.QueryParam("page"), 10, 64)

	if page <= 0 || err != nil {
		page = 1
	}

	data, err := h.service.List(c.Request().Context(), uint(page))

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Employees successfully fetched!",
		Meta: &response.MetaPagination{
			CurrentPage: data.Meta.CurrentPage,
			Limit:       data.Meta.Limit,
			TotalData:   data.Meta.TotalData,
			TotalPages:  data.Meta.TotalPages,
		},
		Data: data.List,
	})
}

func New(service Service, validate *validator.Validate) *Handler {
	return &Handler{
		service:  service,
		validate: validate,
	}
}
