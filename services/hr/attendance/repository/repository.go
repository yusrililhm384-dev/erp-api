package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math"
	"time"

	"enterprise_resource_planning/services/hr/entity"
	"enterprise_resource_planning/services/hr/readmodel"

	errs "enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

const (
	userAttendanceSummaryQuery = `
		select 
			e.id,
			e.name,
			count(case when a.status = 'present' then 1 end) as present,
			count(case when a.status = 'late' then 1 end) as late,
			count(case when a.status = 'sick' then 1 end) as sick,
			count(case when a.status = 'absent' then 1 end) as absent,
			coalesce(
				sum (
					case when cast(a.check_in as time) > '09:30:00'
					then extract(epoch from(cast(a.check_in as time)- time '09:30:00')) / 60
					else 0
				)
			, 0) as late_duration
		from attendances a inner join employees e on a.employee_id = e.id
		where e.id = $1 and a.date >= $2::DATE and a.date < $3::DATE
	`

	userAttendanceListQuery = `
		select 
			a.employee_id as employee_id,
			e.name,
			coalesce(
				jsonb_agg(
					json_build_object(
						'date', a.date,
						'check_in', a.check_in,
						'check_out', a.check_out,
						'status', a.status,
						'late_duration', coalesce(
							sum (
								case when cast(check_in as time) > '09:30:00'
								then extract(epoch from(cast(check_in as time)- time '09:30:00')) / 60
								else 0
							)
						, 0),
						'note', a.note,
						'attachment_key', a.attachment_key
					)
				) filter (where a.date is not null), '[]'::jsonb
			) as attendance_list
		from employees e left join attendances a on a.employee_id = e.id and a.date >= $2::DATE and a.date < $3::DATE
		where e.id = $1 group by e.id, e.name
	`

	usersAttendancesSummaryQuery = `
		select
			(select count (id) from employees where deleted_at is null) as total_employees,
			count(employee_id) as total_attendances,
			count(case when status = 'present' then 1 end) as present,
			count(case when status = 'late' then 1 end) as late,
			count(case when status = 'sick' then 1 end) as sick,
			count(case when status = 'absent' then 1 end) as absent,
		from attendances
		where a.date = $1::date
	`

	usersAttendancesListQuery = `
		select
			a.date,
			e.id,
			e.name,
			a.check_in,
			a.check_out,
			a.status,
			CASE 
				WHEN CAST(a.check_in AS TIME) > '09:30:00'
				THEN EXTRACT(EPOCH FROM (CAST(a.check_in AS TIME) - TIME '09:30:00')) / 60
				ELSE 0
			END AS late_duration,
			a.note,
			a.attachment_key
		from attendances a inner join employees e on a.employee_id = e.id and e.deleted_at is null
		where a.date = $1::date
	`

	countUsersAttendancesListQuery = `
		select count(employee_id) from attendances where date = $1 order by check_in desc limit $2 offset $3
	`
)

type Repository struct {
}

func (r *Repository) UsersAttendancesList(ctx context.Context, pg *sql.DB, date time.Time, page uint) (*readmodel.UsersAttendancesListResponse, error) {
	data := make([]*readmodel.UsersAttendancesList, 0)

	var count uint

	if err := pg.QueryRowContext(ctx, countUsersAttendancesListQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrAttendanceNotFound
		}

		return nil, err
	}

	const limit = 50

	offset := (page - 1) * limit

	rows, err := pg.QueryContext(ctx, userAttendanceListQuery, date, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		list := &readmodel.UsersAttendancesList{}

		if err := rows.Scan(
			&list.Date,
			&list.EmployeeID,
			&list.Name,
			&list.CheckIn,
			&list.CheckOut,
			&list.Status,
			&list.LateDuration,
			&list.Note,
			&list.AttachmentKey,
		); err != nil {
			return nil, err
		}

		data = append(data, list)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.UsersAttendancesListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) UsersAttendancesSummary(ctx context.Context, pg *sql.DB, date time.Time) (*readmodel.UsersAttendancesSummary, error) {
	data := &readmodel.UsersAttendancesSummary{}

	data.CurrentDate = date
	data.Status = &readmodel.AttendanceStatus{}

	if err := pg.QueryRowContext(ctx, usersAttendancesSummaryQuery, date).Scan(
		&data.EmployeesTotal,
		&data.AttendancesTotal,
		&data.Status.Present,
		&data.Status.Late,
		&data.Status.Sick,
		&data.Status.Absent,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrAttendanceNotFound
		}

		return nil, err
	}

	return data, nil
}

func (r *Repository) UserAttendanceList(ctx context.Context, pg *sql.DB, userId uint, startDate, endDate time.Time) (*readmodel.UserAttendanceList, error) {
	list := &readmodel.UserAttendanceList{}

	list.Periode = &readmodel.Periode{
		StartDate: startDate,
		EndDate:   endDate,
	}

	var data []byte

	if err := pg.QueryRowContext(ctx, userAttendanceListQuery, userId, startDate, endDate).Scan(
		&list.UserID,
		&list.Name,
		&data,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}

		return nil, err
	}

	if err := json.Unmarshal(data, &list.List); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *Repository) UserAttendanceSummary(ctx context.Context, pg *sql.DB, userId uint, startDate, endDate time.Time) (*readmodel.UserAttendanceSummary, error) {
	summary := &readmodel.UserAttendanceSummary{}

	summary.Periode = &readmodel.Periode{
		StartDate: startDate,
		EndDate:   endDate,
	}

	summary.Status = &readmodel.AttendanceStatus{}

	if err := pg.QueryRowContext(ctx, userAttendanceSummaryQuery, userId, startDate, endDate).Scan(
		&summary.EmployeeID,
		&summary.Name,
		&summary.Status.Present,
		&summary.Status.Late,
		&summary.Status.Sick,
		&summary.Status.Absent,
		&summary.LateDurationMinutes,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}

		return nil, err
	}

	return summary, nil
}

func (r *Repository) Create(ctx context.Context, pg *gorm.DB, att *entity.Attendance) error {
	res := pg.Create(att)

	if err := res.Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrForeignKeyViolated):
			return errs.ErrUserNotFound
		case errors.Is(err, gorm.ErrDuplicatedKey):
			return errs.ErrDuplicateAttendanceError
		default:
			return err
		}
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, pg *gorm.DB, att *entity.Attendance) error {
	res := pg.Model(&entity.Attendance{}).Where("employee_id = ? and check_out is null", att.EmployeeID).Updates(att)

	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.ErrDuplicateAttendanceError
		}

		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}

func New() *Repository {
	return &Repository{}
}
