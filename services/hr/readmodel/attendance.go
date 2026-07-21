package readmodel

import (
	"time"

	"enterprise_resource_planning/internal/util"
)

type UsersAttendancesListResponse struct {
	Meta *MetaPagination         `json:"meta"`
	List []*UsersAttendancesList `json:"list"`
}

type UsersAttendancesList struct {
	Date          util.JSONDate `json:"date"`
	EmployeeID    uint          `json:"employee_id"`
	Name          string        `json:"name"`
	CheckIn       *time.Time    `json:"check_in"`
	CheckOut      *time.Time    `json:"check_out"`
	Status        string        `json:"status"`
	LateDuration  uint          `json:"late_duration"`
	Note          *string       `json:"note,omitempty"`
	AttachmentKey *string       `json:"attachment_key,omitempty"`
}

type UsersAttendancesSummary struct {
	EmployeesTotal   uint              `json:"employees_total"`
	AttendancesTotal uint              `json:"attendances_total"`
	CurrentDate      time.Time         `json:"current_date"`
	Status           *AttendanceStatus `json:"status"`
}

type UserAttendanceSummary struct {
	EmployeeID          uint              `json:"employee_id"`
	Name                string            `json:"name"`
	Periode             *Periode          `json:"periode"`
	Status              *AttendanceStatus `json:"status"`
	LateDurationMinutes uint              `json:"late_duration_minutes"`
}

type AttendanceStatus struct {
	Present uint `json:"present"`
	Late    uint `json:"late"`
	Sick    uint `json:"sick"`
	Absent  uint `json:"absent"`
}

type Periode struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type UserAttendanceList struct {
	UserID  uint              `json:"user_id"`
	Name    string            `json:"name"`
	Periode *Periode          `json:"periode"`
	List    []*AttendanceList `json:"attendances"`
}

type AttendanceList struct {
	Date          util.JSONDate `json:"date"`
	CheckIn       *time.Time    `json:"check_in"`
	CheckOut      *time.Time    `json:"check_out"`
	Status        string        `json:"status"`
	LateDuration  uint          `json:"late_duration"`
	Note          *string       `json:"note,omitempty"`
	AttachmentKey *string       `json:"attachment_key,omitempty"`
}
