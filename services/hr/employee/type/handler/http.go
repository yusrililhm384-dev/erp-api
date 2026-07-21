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
	Create(ctx context.Context, userId uint, req *dto.CreateTypeReq) error
	Delete(ctx context.Context, userId, typeId uint) error
	Update(ctx context.Context, userId uint, typeId uint, req *dto.UpdateTypeReq) error
	Detail(ctx context.Context, typeId uint) (*dto.DetailTypeRes, error)
	List(ctx context.Context) ([]*readmodel.Type, error)
}

type Handler struct {
	service  Service
	validate *validator.Validate
}

func (h *Handler) Create(c *echo.Context) error {
	c.Logger().Info("Create Type")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	req := &dto.CreateTypeReq{}

	if err := c.Bind(req); err != nil {
		c.Logger().Warn(err.Error())

		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid body request")
	}

	if err := h.validate.Struct(req); err != nil {
		c.Logger().Warn(err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, "Missing required body request")
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
		Message: "Type successfully created!",
	})
}

func (h *Handler) Update(c *echo.Context) error {
	c.Logger().Info("Update Type")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	typeId, err := strconv.ParseUint(c.Param("typeId"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid type id!")
	}

	req := &dto.UpdateTypeReq{}

	if err := c.Bind(req); err != nil {
		c.Logger().Warn(err.Error())

		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid body request")
	}

	if err := h.validate.Struct(req); err != nil {
		c.Logger().Warn(err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, "Missing required body request")
	}

	if err := h.service.Update(c.Request().Context(), sessionData.UserId, uint(typeId), req); err != nil {
		c.Logger().Warn(err.Error())

		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "User not found!")
		case errors.Is(err, errs.ErrTypeNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Type not found!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Type successfully updated!",
	})
}

func (h *Handler) Delete(c *echo.Context) error {
	c.Logger().Info("Delete Type")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	typeId, err := strconv.ParseUint(c.Param("typeId"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid type id!")
	}

	if err := h.service.Delete(c.Request().Context(), sessionData.UserId, uint(typeId)); err != nil {
		c.Logger().Warn(err.Error())

		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "User not found!")
		case errors.Is(err, errs.ErrTypeNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Type not found!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Type successfully deleted!",
	})
}

func (h *Handler) Detail(c *echo.Context) error {
	c.Logger().Info("Type Detail")

	typeId, err := strconv.ParseUint(c.Param("typeId"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid type id!")
	}

	data, err := h.service.Detail(c.Request().Context(), uint(typeId))

	if err != nil {
		c.Logger().Warn(err.Error())

		if errors.Is(err, errs.ErrTypeNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Type not found!")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Type successfully fetched!",
		Data:    data,
	})
}

func (h *Handler) List(c *echo.Context) error {
	c.Logger().Info("Type List")

	data, err := h.service.List(c.Request().Context())

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Types successfully fetched!",
		Data:    data,
	})
}

func New(service Service, validate *validator.Validate) *Handler {
	return &Handler{
		service:  service,
		validate: validate,
	}
}
