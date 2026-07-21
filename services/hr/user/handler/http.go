package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"enterprise_resource_planning/internal/shared/response"
	"enterprise_resource_planning/services/hr/dto"

	errs "enterprise_resource_planning/internal/shared/errors"
	auth "enterprise_resource_planning/services/auth/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type Service interface {
	ChangePassword(ctx context.Context, userId uint, req *dto.ChangePasswordReq) error
	ChangeUsername(ctx context.Context, userId uint, req *dto.ChangeUsernameReq) error
	Login(ctx context.Context, req *dto.UserLoginReq) (*dto.UserLoginResponse, error)
	Logout(ctx context.Context, sessionId string) error
}

type Handler struct {
	service  Service
	validate *validator.Validate
}

func (h *Handler) Login(c *echo.Context) error {
	c.Logger().Info("User Login")

	req := &dto.UserLoginReq{}

	if err := c.Bind(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid nody request!")
	}

	if err := h.validate.Struct(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields!")
	}

	res, err := h.service.Login(c.Request().Context(), req)

	if err != nil {
		c.Logger().Warn(err.Error())

		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "User not found!")
		case errors.Is(err, errs.ErrTooManyRequest):
			return echo.NewHTTPError(http.StatusTooManyRequests, "Too many request your account will be locked!")
		case errors.Is(err, errs.ErrUserLocked):
			return echo.NewHTTPError(http.StatusTooManyRequests, "Your account is locked!")
		case errors.Is(err, errs.ErrInvalidPassword):
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid password!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	c.SetCookie(&http.Cookie{
		Name:     "session_id",
		Value:    res.SessionId,
		Path:     "/",
		Expires:  time.Now().Add(8 * time.Hour),
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "User successfully logged in!",
	})
}

func (h *Handler) ChangePassword(c *echo.Context) error {
	c.Logger().Info("User Change Password")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	req := &dto.ChangePasswordReq{}

	if err := c.Bind(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid nody request!")
	}

	if err := h.validate.Struct(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields!")
	}

	if err := h.service.ChangePassword(c.Request().Context(), sessionData.UserId, req); err != nil {
		c.Logger().Warn(err.Error())

		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "User not found!")
		case errors.Is(err, errs.ErrInvalidPassword):
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid password!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	c.SetCookie(&http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		MaxAge:   -1,
	})

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Password successfully updated!",
	})
}

func (h *Handler) ChangeUsername(c *echo.Context) error {
	c.Logger().Info("Change Username")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	req := &dto.ChangeUsernameReq{}

	if err := c.Bind(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid nody request!")
	}

	if err := h.validate.Struct(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields!")
	}

	if err := h.service.ChangeUsername(c.Request().Context(), sessionData.UserId, req); err != nil {
		c.Logger().Warn(err.Error())

		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "User not found!")
		case errors.Is(err, errs.ErrDuplicateUsernameError):
			return echo.NewHTTPError(http.StatusConflict, "Username already used!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Username successfully updated!",
	})
}

func (h *Handler) Logout(c *echo.Context) error {
	c.Logger().Info("User Logout")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	if err := h.service.Logout(c.Request().Context(), sessionData.SessionId); err != nil {
		c.Logger().Warn(err.Error())

		if errors.Is(err, errs.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found!")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	c.SetCookie(&http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
	})

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "User successfully logged out!",
	})
}

func New(service Service, validate *validator.Validate) *Handler {
	return &Handler{
		service:  service,
		validate: validate,
	}
}
