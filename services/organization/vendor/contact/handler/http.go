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
	Create(ctx context.Context, userId uint, vendorId uint, req *dto.ContactVendorReq) error
	Update(ctx context.Context, userId uint, vendorId uint, contactId uint, req *dto.UpdateContactVendorReq) error
	Delete(ctx context.Context, userId uint, vendorId uint, contactId uint) error
	Detail(ctx context.Context, contactId, vendorId uint) (*readmodel.VendorContactDetail, error)
}

type Handler struct {
	service  Service
	validate *validator.Validate
}

func (h *Handler) Create(c *echo.Context) error {
	c.Logger().Info("Create Contact")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	vendorId, err := strconv.ParseUint(c.Param("vendorId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid vendor id")
	}

	req := &dto.ContactVendorReq{}

	if err := c.Bind(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid body request!")
	}

	if err := h.validate.Struct(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields!")
	}

	if err := h.service.Create(c.Request().Context(), sessionData.UserId, uint(vendorId), req); err != nil {
		c.Logger().Warn(err.Error())

		switch {
		case errors.Is(err, errs.ErrDuplicateContactError):
			return echo.NewHTTPError(http.StatusConflict, "Contact already exists!")
		case errors.Is(err, errs.ErrForeignKeyError):
			return echo.NewHTTPError(http.StatusNotFound, "Invalid reference id!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	return c.JSON(http.StatusCreated, &response.Response{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Contact successfully created!",
	})
}

func (h *Handler) Update(c *echo.Context) error {
	c.Logger().Info("Update Contact")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	vendorId, err := strconv.ParseUint(c.Param("vendorId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid vendor id")
	}

	contactId, err := strconv.ParseUint(c.Param("contactId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid contact id")
	}

	req := &dto.UpdateContactVendorReq{}

	if err := c.Bind(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid body request!")
	}

	if err := h.validate.Struct(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields!")
	}

	if err := h.service.Update(c.Request().Context(), sessionData.UserId, uint(vendorId), uint(contactId), req); err != nil {
		c.Logger().Warn(err.Error())

		switch {
		case errors.Is(err, errs.ErrDuplicateContactError):
			return echo.NewHTTPError(http.StatusConflict, "Contact already exists!")
		case errors.Is(err, errs.ErrForeignKeyError):
			return echo.NewHTTPError(http.StatusNotFound, "Invalid reference id!")
		case errors.Is(err, errs.ErrVendorContactNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Vendor contact not found!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Contact successfully updated!",
	})
}

func (h *Handler) Delete(c *echo.Context) error {
	c.Logger().Info("Delete Contact")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	vendorId, err := strconv.ParseUint(c.Param("vendorId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid vendor id")
	}

	contactId, err := strconv.ParseUint(c.Param("contactId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid contact id")
	}

	if err := h.service.Delete(c.Request().Context(), sessionData.UserId, uint(vendorId), uint(contactId)); err != nil {
		c.Logger().Warn(err.Error())

		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "User not found!")
		case errors.Is(err, errs.ErrVendorContactNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Vendor contact not found!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Contact successfully deleted!",
	})
}

func (h *Handler) Detail(c *echo.Context) error {
	c.Logger().Info("Vendor Contact Detail")

	vendorId, err := strconv.ParseUint(c.Param("vendorId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid vendor id")
	}

	contactId, err := strconv.ParseUint(c.Param("contactId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid contact id")
	}

	data, err := h.service.Detail(c.Request().Context(), uint(contactId), uint(vendorId))

	if err != nil {
		c.Logger().Warn(err.Error())

		if errors.Is(err, errs.ErrVendorContactNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Vendor contact not found!")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Contact successfully fetched!",
		Data:    data,
	})
}

func New(service Service, validate *validator.Validate) *Handler {
	return &Handler{
		service:  service,
		validate: validate,
	}
}
