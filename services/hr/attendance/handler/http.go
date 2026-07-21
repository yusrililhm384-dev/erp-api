package handler

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"enterprise_resource_planning/internal/shared/response"
	"enterprise_resource_planning/services/hr/dto"
	"enterprise_resource_planning/services/hr/readmodel"

	errs "enterprise_resource_planning/internal/shared/errors"
	auth "enterprise_resource_planning/services/auth/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type AttendanceService interface {
	CheckIn(ctx context.Context, userId uint) error
	Sick(ctx context.Context, userId uint, req *dto.CreateSickReq) error
	CheckOut(ctx context.Context, userId uint) error
	UserAttendanceList(ctx context.Context, userId uint, req *dto.FilterDateReq) (*readmodel.UserAttendanceList, error)
	UserAttendanceSummary(ctx context.Context, userId uint, req *dto.FilterDateReq) (*readmodel.UserAttendanceSummary, error)
	UsersAttendancesSummary(ctx context.Context, date time.Time) (*readmodel.UsersAttendancesSummary, error)
	UsersAttendancesList(ctx context.Context, date time.Time, page uint) (*readmodel.UsersAttendancesListResponse, error)
}

type StorageService interface {
	Upload(ctx context.Context, objKey string, file *multipart.FileHeader) error
}

type Handler struct {
	attendanceService AttendanceService
	storageService    StorageService
	validate          *validator.Validate
}

func (h *Handler) UsersAttendancesList(c *echo.Context) error {
	c.Logger().Info("Users Attendances Summary")

	dateString := c.QueryParam("date")

	loc, err := time.LoadLocation("Asia/Jakarta")

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	var date time.Time

	if dateString == "" {
		date = time.Now().In(loc)
	} else {
		date, err = time.ParseInLocation(time.DateOnly, dateString, loc)

		if err != nil {
			c.Logger().Warn(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}

	}

	page, err := strconv.ParseUint(c.QueryParam("page"), 10, 64)

	if page <= 0 || err != nil {
		page = 1
	}

	data, err := h.attendanceService.UsersAttendancesList(c.Request().Context(), date, uint(page))

	if err != nil {
		c.Logger().Warn(err.Error())

		if errors.Is(err, errs.ErrAttendanceNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Users attendances not found!")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Users attendances list successfully fetched!",
		Meta: &response.MetaPagination{
			CurrentPage: data.Meta.CurrentPage,
			Limit:       data.Meta.Limit,
			TotalData:   data.Meta.TotalData,
			TotalPages:  data.Meta.TotalPages,
		},
		Data: data.List,
	})
}

func (h *Handler) UsersAttendancesSummary(c *echo.Context) error {
	c.Logger().Info("Users Attendances Summary")

	dateString := c.QueryParam("date")

	loc, err := time.LoadLocation("Asia/Jakarta")

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	var date time.Time

	if dateString == "" {
		date = time.Now().In(loc)
	} else {
		date, err = time.ParseInLocation(time.DateOnly, dateString, loc)

		if err != nil {
			c.Logger().Warn(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}

	}

	data, err := h.attendanceService.UsersAttendancesSummary(c.Request().Context(), date)

	if err != nil {
		c.Logger().Warn(err.Error())

		if errors.Is(err, errs.ErrAttendanceNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Users attendances summary not found!")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Users attendances summary successfully fetched!",
		Data:    data,
	})
}

func (h *Handler) UserDetailSummary(c *echo.Context) error {
	c.Logger().Info("User Detail Summary")

	userId, err := strconv.ParseUint(c.Param("userId"), 10, 64)

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user id!")
	}

	req := &dto.FilterDateStringReq{}

	if err := c.Bind(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid filter date!")
	}

	if err := h.validate.Struct(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields!")
	}

	date, err := req.DateStringToDate()

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	data, err := h.attendanceService.UserAttendanceSummary(c.Request().Context(), uint(userId), date)

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
		Message: "User attendance summary successfully fetched!",
		Data:    data,
	})
}

func (h *Handler) UserAttendanceSummary(c *echo.Context) error {
	c.Logger().Info("User Attendance Summary")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	req := &dto.FilterDateStringReq{}

	if err := c.Bind(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid filter date!")
	}

	if err := h.validate.Struct(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields!")
	}

	date, err := req.DateStringToDate()

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	data, err := h.attendanceService.UserAttendanceSummary(c.Request().Context(), sessionData.UserId, date)

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
		Message: "User attendance summary successfully fetched!",
		Data:    data,
	})
}

func (h *Handler) UserAttendanceList(c *echo.Context) error {
	c.Logger().Info("User Attendance List")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	req := &dto.FilterDateStringReq{}

	if err := c.Bind(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid filter date!")
	}

	if err := h.validate.Struct(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields!")
	}

	date, err := req.DateStringToDate()

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	data, err := h.attendanceService.UserAttendanceList(c.Request().Context(), sessionData.UserId, date)

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
		Message: "User attendance list successfully fetched!",
		Data:    data,
	})
}

func (h *Handler) CheckIn(c *echo.Context) error {
	c.Logger().Info("User Attendance")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	if err := h.attendanceService.CheckIn(c.Request().Context(), sessionData.UserId); err != nil {
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "User not found!")
		case errors.Is(err, errs.ErrTodayIsWeekend):
			return echo.NewHTTPError(http.StatusBadRequest, "Today is weekend! Do You wanna overtime?")
		case errors.Is(err, errs.ErrDuplicateAttendanceError):
			return echo.NewHTTPError(http.StatusConflict, "Duplicate attendance!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	return c.JSON(http.StatusCreated, &response.Response{
		Success: true,
		Status:  http.StatusCreated,
		Message: "User successfully attendanced!",
	})
}

func (h *Handler) Sick(c *echo.Context) error {
	c.Logger().Info("User Sick")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	req := &dto.SickReq{}

	if err := c.Bind(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid filter date!")
	}

	if err := h.validate.Struct(req); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields!")
	}

	file, err := c.FormFile("file")

	if err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to upload file!")
	}

	const maxSize = 500 * 1024

	if file.Size > maxSize {
		return echo.NewHTTPError(http.StatusBadRequest, "File too big!")
	}

	ext := filepath.Ext(file.Filename)

	if ext != ".pdf" && ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file type!")
	}

	objectKey := fmt.Sprintf("sick_letter/user_%d/%d%s", sessionData.UserId, time.Now().Unix(), ext)

	if err := h.storageService.Upload(c.Request().Context(), objectKey, file); err != nil {
		c.Logger().Warn(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Internal server error")
	}

	if err := h.attendanceService.Sick(c.Request().Context(), sessionData.UserId, &dto.CreateSickReq{
		Note:          req.Note,
		AttachmentKey: objectKey,
	}); err != nil {
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "User not found!")
		case errors.Is(err, errs.ErrTodayIsWeekend):
			return echo.NewHTTPError(http.StatusBadRequest, "Today is weekend! Do You wanna overtime?")
		case errors.Is(err, errs.ErrDuplicateAttendanceError):
			return echo.NewHTTPError(http.StatusConflict, "Duplicate attendance!")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
		}
	}

	return c.JSON(http.StatusCreated, &response.Response{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Sick successfully requested!",
	})
}

func (h *Handler) CheckOut(c *echo.Context) error {
	c.Logger().Info("User Check Out")

	sessionData, err := auth.GetSessionFromContext(c)

	if err != nil {
		c.Logger().Warn(err.Error())
		return err
	}

	if err := h.attendanceService.CheckOut(c.Request().Context(), sessionData.UserId); err != nil {
		c.Logger().Warn(err.Error())

		if errors.Is(err, errs.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found!")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Oops! Something went wrong!")
	}

	return c.JSON(http.StatusOK, &response.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "User successfully checkout!",
	})
}

func New(attendanceService AttendanceService, storageService StorageService, validate *validator.Validate) *Handler {
	return &Handler{
		attendanceService: attendanceService,
		storageService:    storageService,
		validate:          validate,
	}
}
