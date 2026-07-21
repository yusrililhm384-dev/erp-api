package service

import (
	"context"
	"database/sql"
	"time"

	"enterprise_resource_planning/internal/util"
	"enterprise_resource_planning/services/hr/dto"
	"enterprise_resource_planning/services/hr/entity"
	"enterprise_resource_planning/services/hr/readmodel"

	errs "enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, pg *gorm.DB, att *entity.Attendance) error
	Update(ctx context.Context, pg *gorm.DB, att *entity.Attendance) error
	UserAttendanceList(ctx context.Context, pg *sql.DB, userId uint, startDate, endDate time.Time) (*readmodel.UserAttendanceList, error)
	UserAttendanceSummary(ctx context.Context, pg *sql.DB, userId uint, startDate, endDate time.Time) (*readmodel.UserAttendanceSummary, error)
	UsersAttendancesSummary(ctx context.Context, pg *sql.DB, date time.Time) (*readmodel.UsersAttendancesSummary, error)
	UsersAttendancesList(ctx context.Context, pg *sql.DB, date time.Time, page uint) (*readmodel.UsersAttendancesListResponse, error)
}

type Service struct {
	pg         *gorm.DB
	repository Repository
}

func (s *Service) UsersAttendancesList(ctx context.Context, date time.Time, page uint) (*readmodel.UsersAttendancesListResponse, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	data, err := s.repository.UsersAttendancesList(ctx, pg, date, page)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Service) UsersAttendancesSummary(ctx context.Context, date time.Time) (*readmodel.UsersAttendancesSummary, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	data, err := s.repository.UsersAttendancesSummary(ctx, pg, date)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Service) UserAttendanceSummary(ctx context.Context, userId uint, req *dto.FilterDateReq) (*readmodel.UserAttendanceSummary, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	data, err := s.repository.UserAttendanceSummary(ctx, pg, userId, req.StartDate, req.EndDate)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Service) UserAttendanceList(ctx context.Context, userId uint, req *dto.FilterDateReq) (*readmodel.UserAttendanceList, error) {
	pg, err := s.pg.DB()

	if err != nil {
		return nil, err
	}

	data, err := s.repository.UserAttendanceList(ctx, pg, userId, req.StartDate, req.EndDate)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Service) CheckIn(ctx context.Context, userId uint) error {
	loc := time.FixedZone("WIB", 7*60*60)

	now := time.Now().In(loc)

	if err := s.checkIsWeekend(now); err != nil {
		return err
	}

	att := &entity.Attendance{
		Date: util.JSONDate{
			Time: now,
		},
		CheckIn:    &now,
		Status:     "present",
		EmployeeID: userId,
	}

	if now.Hour() >= 9 {
		att.Status = "late"
	}

	if err := s.repository.Create(ctx, s.pg, att); err != nil {
		return err
	}

	return nil
}

func (s *Service) Sick(ctx context.Context, userId uint, req *dto.CreateSickReq) error {
	loc := time.FixedZone("WIB", 7*60*60)

	now := time.Now().In(loc)

	if err := s.checkIsWeekend(now); err != nil {
		return err
	}

	if err := s.checkIsWeekend(now); err != nil {
		return err
	}

	att := &entity.Attendance{
		Date: util.JSONDate{
			Time: now,
		},
		Status:        "sick",
		EmployeeID:    userId,
		Note:          &req.Note,
		AttachmentKey: &req.AttachmentKey,
	}

	if err := s.repository.Create(ctx, s.pg, att); err != nil {
		return err
	}

	return nil
}

func (s *Service) CheckOut(ctx context.Context, userId uint) error {
	loc := time.FixedZone("WIB", 7*60*60)

	now := time.Now().In(loc)

	if err := s.checkIsWeekend(now); err != nil {
		return err
	}

	att := &entity.Attendance{
		CheckOut:   &now,
		EmployeeID: userId,
	}

	if err := s.repository.Update(ctx, s.pg, att); err != nil {
		return err
	}

	return nil
}

func (s *Service) checkIsWeekend(now time.Time) error {
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return errs.ErrTodayIsWeekend
	}

	return nil
}

func New(pg *gorm.DB, repository Repository) *Service {
	return &Service{
		pg:         pg,
		repository: repository,
	}
}
