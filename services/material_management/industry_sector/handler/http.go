package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"enterprise_resource_planning/internal/shared/response"
	"enterprise_resource_planning/services/material_management/dto"
	"enterprise_resource_planning/services/material_management/readmodel"

	errs "enterprise_resource_planning/internal/shared/errors"
	auth "enterprise_resource_planning/services/auth/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type Service interface {
	List(ctx context.Context, page uint) (*readmodel.IndustrySectorListResponse, error)
	Detail(ctx context.Context, industrySectorId uint) (*readmodel.IndustrySectorDetail, error)
	Create(ctx context.Context, userId uint, req *dto.CreateIndustrySectorReq) error
	Update(ctx context.Context, userId uint, industrySectorId uint, req *dto.UpdateIndustrySectorReq) error
	Delete(ctx context.Context, userId uint, industrySectorId uint) error
}

type Handler struct {
	service  Service
	validate *validator.Validate
}

func (h *Handler) List(c *echo.Context) error {
	c.Logger().Info("Industry Sector List")

	page, err := strconv.ParseUint(c.QueryParam("page"), 10, 64)

	if page <= 0 || err != nil {
		page = 1
	}

	data, err := h.service.List(c.Request().Context(), uint(page))

	if err != nil {
		c.Logger().Warn(err.Error())

		if errors.Is(err, errs.ErrIndustrySectorNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Industry sector not found!")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Industry sector list successfully fetched!",
		Meta: &response.MetaPagination{
			CurrentPage: data.Meta.CurrentPage,
			Limit:       data.Meta.Limit,
			TotalData:   data.Meta.TotalData,
			TotalPages:  data.Meta.TotalPages,
		},
		Data: data.List,
	})
}

func (h *Handler) Detail(c *echo.Context) error {
	c.Logger().Info("Industry Sector Detail")

	industrySectorId, err := strconv.ParseUint(c.Param("industrySectorId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid material group id")
	}

	data, err := h.service.Detail(c.Request().Context(), uint(industrySectorId))

	if err != nil {
		c.Logger().Warn(err.Error())

		if errors.Is(err, errs.ErrIndustrySectorNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Industry sector not found!")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Industry  sector detail successfully fetched!",
		Data:    data,
	})
}

func (h *Handler) Create(c *echo.Context) error {
	c.Logger().Info("Create Industry Sector")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	req := &dto.CreateIndustrySectorReq{}

	if err := c.Bind(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid body request!")
	}

	if err := h.validate.Struct(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields!")
	}

	if err := h.service.Create(c.Request().Context(), sessionData.UserId, req); err != nil {
		c.Logger().Warn(err.Error())

		if errors.Is(err, errs.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found!")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	return c.JSON(http.StatusCreated, &response.Response{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Industry sector successfully created!",
	})
}

func (h *Handler) Update(c *echo.Context) error {
	c.Logger().Info("Update Industry Sector")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	industrySectorId, err := strconv.ParseUint(c.Param("industrySectorId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid material group id")
	}

	req := &dto.UpdateIndustrySectorReq{}

	if err := c.Bind(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid body request!")
	}

	if err := h.validate.Struct(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields!")
	}

	if err := h.service.Update(c.Request().Context(), sessionData.UserId, uint(industrySectorId), req); err != nil {
		c.Logger().Warn(err.Error())

		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "User not found!")
		case errors.Is(err, errs.ErrIndustrySectorNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Industry sector not found!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Industry sector successfully updated!",
	})
}

func (h *Handler) Delete(c *echo.Context) error {
	c.Logger().Info("Delete Industry Sector")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	industrySectorId, err := strconv.ParseUint(c.Param("industrySectorId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid material group id")
	}

	if err := h.service.Delete(c.Request().Context(), sessionData.UserId, uint(industrySectorId)); err != nil {
		c.Logger().Warn(err.Error())

		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "User not found!")
		case errors.Is(err, errs.ErrIndustrySectorNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Industry sector not found!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Industry sector successfully deleted!",
	})
}

func New(service Service, validate *validator.Validate) *Handler {
	return &Handler{
		service:  service,
		validate: validate,
	}
}
