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
	List(ctx context.Context, page uint) (*readmodel.PaymentMethodListResponse, error)
	Detail(ctx context.Context, paymentMethodId uint) (*readmodel.PaymentMethodDetail, error)
	Create(ctx context.Context, userId uint, req *dto.CreatePaymentMethodReq) error
	Update(ctx context.Context, userId uint, paymentMethodId uint, req *dto.UpdatePaymentMethodReq) error
	Delete(ctx context.Context, userId uint, paymentMethodId uint) error
}

type Handler struct {
	service  Service
	validate *validator.Validate
}

func (h *Handler) List(c *echo.Context) error {
	c.Logger().Info("Payment Method List")

	page, err := strconv.ParseUint(c.QueryParam("page"), 10, 64)

	if page <= 0 || err != nil {
		page = 1
	}

	data, err := h.service.List(c.Request().Context(), uint(page))

	if err != nil {
		c.Logger().Warn(err.Error())

		if errors.Is(err, errs.ErrPaymentMethodNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Payment method not found!")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Payment method list successfully fetched!",
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
	c.Logger().Info("Payment Method Detail")

	paymentMethodId, err := strconv.ParseUint(c.Param("paymentMethodId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payment method id!")
	}

	data, err := h.service.Detail(c.Request().Context(), uint(paymentMethodId))

	if err != nil {
		c.Logger().Warn(err.Error())

		if errors.Is(err, errs.ErrPaymentMethodNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Payment method not found!")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Payment method detail successfully fetched!",
		Data:    data,
	})
}

func (h *Handler) Create(c *echo.Context) error {
	c.Logger().Info("Create Payment Method")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	req := &dto.CreatePaymentMethodReq{}

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

		if errors.Is(err, errs.ErrDuplicateError) {
			return echo.NewHTTPError(http.StatusConflict, "Payment term already exists!")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	return c.JSON(http.StatusCreated, &response.Response{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Payment method successfully created!",
	})
}

func (h *Handler) Update(c *echo.Context) error {
	c.Logger().Info("Update Payment Method")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	paymentMethodId, err := strconv.ParseUint(c.Param("paymentMethodId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payment method id!")
	}

	req := &dto.UpdatePaymentMethodReq{}

	if err := c.Bind(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid body request!")
	}

	if err := h.validate.Struct(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields!")
	}

	if err := h.service.Update(c.Request().Context(), sessionData.UserId, uint(paymentMethodId), req); err != nil {
		c.Logger().Warn(err.Error())

		switch {
		case errors.Is(err, errs.ErrForeignKeyError):
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid reference id!")
		case errors.Is(err, errs.ErrPaymentMethodNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Payment method not found!")
		case errors.Is(err, errs.ErrDuplicateError):
			return echo.NewHTTPError(http.StatusConflict, "Payment term already exists!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Payment method successfully updated!",
	})
}

func (h *Handler) Delete(c *echo.Context) error {
	c.Logger().Info("Delete Payment Method")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	paymentMethodId, err := strconv.ParseUint(c.Param("paymentMethodId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payment method id!")
	}

	if err := h.service.Delete(c.Request().Context(), sessionData.UserId, uint(paymentMethodId)); err != nil {
		c.Logger().Warn(err.Error())

		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "User not found!")
		case errors.Is(err, errs.ErrPaymentMethodNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Payment method not found!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Payment method successfully deleted!",
	})
}

func New(service Service, validate *validator.Validate) *Handler {
	return &Handler{
		service:  service,
		validate: validate,
	}
}
