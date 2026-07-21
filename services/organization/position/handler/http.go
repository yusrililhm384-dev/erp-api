package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"enterprise_resource_planning/internal/shared/response"
	"enterprise_resource_planning/services/organization/dto"
	"enterprise_resource_planning/services/organization/readmodel"

	errs "enterprise_resource_planning/internal/shared/errors"
	auth "enterprise_resource_planning/services/auth/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type Service interface {
	List(ctx context.Context, page uint) (*readmodel.PositionListResponse, error)
	Detail(ctx context.Context, positionId uint) (*readmodel.PositionDetail, error)
	Create(ctx context.Context, userId uint, req *dto.CreatePositionReq) error
	Update(ctx context.Context, userId uint, positionId uint, req *dto.UpdatePositionReq) error
	Delete(ctx context.Context, userId uint, positionId uint) error
}

type Handler struct {
	service  Service
	validate *validator.Validate
}

func (h *Handler) List(c *echo.Context) error {
	c.Logger().Info("Position List")

	page, err := strconv.ParseUint(c.QueryParam("page"), 10, 64)

	if page <= 0 || err != nil {
		page = 1
	}

	data, err := h.service.List(c.Request().Context(), uint(page))

	if err != nil {
		c.Logger().Warn(err.Error())

		if errors.Is(err, errs.ErrPositionNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Position not found!")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Position list successfully fetched!",
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
	c.Logger().Info("Position Detail")

	positionId, err := strconv.ParseUint(c.Param("positionId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid position id")
	}

	data, err := h.service.Detail(c.Request().Context(), uint(positionId))

	if err != nil {
		c.Logger().Warn(err.Error())

		if errors.Is(err, errs.ErrPositionNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Position not found!")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Position detail successfully fetched!",
		Data:    data,
	})
}

func (h *Handler) Create(c *echo.Context) error {
	c.Logger().Info("Create Position")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	req := &dto.CreatePositionReq{}

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

		if errors.Is(err, errs.ErrForeignKeyError) {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid reference id!")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	return c.JSON(http.StatusCreated, &response.Response{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Position successfully created!",
	})
}

func (h *Handler) Update(c *echo.Context) error {
	c.Logger().Info("Update Position")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	positionId, err := strconv.ParseUint(c.Param("positionId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid position id")
	}

	req := &dto.UpdatePositionReq{}

	if err := c.Bind(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid body request!")
	}

	if err := h.validate.Struct(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields!")
	}

	if err := h.service.Update(c.Request().Context(), sessionData.UserId, uint(positionId), req); err != nil {
		c.Logger().Warn(err.Error())

		switch {
		case errors.Is(err, errs.ErrForeignKeyError):
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid reference id!")
		case errors.Is(err, errs.ErrPositionNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Position not found!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Position successfully updated!",
	})
}

func (h *Handler) Delete(c *echo.Context) error {
	c.Logger().Info("Delete Position")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	positionId, err := strconv.ParseUint(c.Param("positionId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid position id")
	}

	if err := h.service.Delete(c.Request().Context(), sessionData.UserId, uint(positionId)); err != nil {
		c.Logger().Warn(err.Error())

		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "User not found!")
		case errors.Is(err, errs.ErrPositionNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Position not found!")
		case errors.Is(err, errs.ErrPositionHasActiveEmployees):
			return echo.NewHTTPError(http.StatusNotFound, "Position has active employees!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Position successfully deleted!",
	})
}

func New(service Service, validate *validator.Validate) *Handler {
	return &Handler{
		service:  service,
		validate: validate,
	}
}
