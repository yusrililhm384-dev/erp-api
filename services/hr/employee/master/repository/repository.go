package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math"

	"enterprise_resource_planning/services/hr/entity"
	"enterprise_resource_planning/services/hr/readmodel"

	errs "enterprise_resource_planning/internal/shared/errors"

	"gorm.io/gorm"
)

const (
	employeeListQuery = `
		SELECT 
			e.id, 
			e.name, 
			e.email, 
			e.phone, 
			e.gender,
			e.join_date, 
			p.name AS position,
			e.created_at,
			e.updated_at
		FROM employees e 
		INNER JOIN mutations m ON m.employee_id = e.id 
		INNER JOIN positions p ON m.position_id = p.id 
		WHERE e.deleted_at IS NULL 
			AND m.deleted_at IS NULL 
			AND p.deleted_at IS NULL
			AND m.end_date IS NULL
		group by created_at desc
		limit $1 offset $2
	`

	countEmployeeQuery = `
		select count(id) from employees where deleted_at is null
	`

	employeeDetailQuery = `
		select
			e.id,
			e.name,
			e.email,
			e.phone,
			e.gender,
			e.address,
			e.join_date,
			l.identity_number,
			l.tax_number,
			l.health_insurance_number,
			l.labor_insurance_number,
			coalesce(
				jsonb_agg(
					json_build_object(
						'position', p.name,
						'branch', b.name,
						'department', d.name,
						'type', t.name,
						'start_date', m.start_date,
						'end_date', m.end_date
					) order by m.end_date nulls first, m.start_date desc
				) filter (where p.id is not null ), '[]'::jsonb
			) as histories,
			e.created_at,
			e.updated_at
		from employees e
		inner join legals l on l.employee_id = e.id and l.deleted_at is null
		left join mutations m on m.employee_id = e.id and m.deleted_at is null
		left join positions p on m.position_id = p.id and p.deleted_at is null
		left join types t on m.type_id = t.id and t.deleted_at is null
		left join departments d on p.department_id = d.id and d.deleted_at is null
		left join branches b on m.branch_id = b.id and b.deleted_at is null
		where e.iGROUPd = $1 and e.deleted_at is null
		group bY
			e.id,
			e.name,
			e.email,
			e.phone,
			e.gender,
			e.address,
			e.join_date,
			l.identity_number,
			l.tax_number,
			l.health_insurance_number,
			l.labor_insurance_number,
			e.created_at,
			e.updated_at;
	`
)

type Repository struct {
}

func (r *Repository) Create(tx *gorm.DB, employee *entity.Employee) error {
	result := tx.Create(employee)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.ErrDuplicateContactError
		}

		return err
	}

	return nil
}

func (r *Repository) Update(tx *gorm.DB, employee *entity.Employee) error {
	result := tx.Model(&entity.Employee{}).Updates(employee)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrDuplicateContactError
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrEmployeeMaserNotFound
	}

	return nil
}

func (r *Repository) Delete(tx *gorm.DB, employeeId uint) error {
	result := tx.Model(&entity.Employee{}).Delete(employeeId)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errs.ErrUserNotFound
		}

		return err
	}

	if result.RowsAffected == 0 {
		return errs.ErrEmployeeMaserNotFound
	}

	return nil
}

func (r *Repository) List(ctx context.Context, pg *sql.DB, page uint) (*readmodel.EmployeeListResponse, error) {
	data := make([]*readmodel.EmployeeList, 0)

	var count uint

	if err := pg.QueryRowContext(ctx, countEmployeeQuery).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrEmployeeNotFound
		}

		return nil, err
	}

	const limit = 50

	offset := (page - 1) * limit

	rows, err := pg.QueryContext(ctx, employeeListQuery, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		employee := &readmodel.EmployeeList{}

		if err := rows.Scan(
			&employee.Id,
			&employee.Name,
			&employee.Email,
			&employee.Phone,
			&employee.Gender,
			&employee.JoinDate,
			&employee.Position,
			&employee.CreatedAt,
			&employee.UpdatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, employee)
	}

	totalPages := uint(math.Ceil(float64(count) / float64(limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &readmodel.EmployeeListResponse{
		Meta: &readmodel.MetaPagination{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   uint(count),
			TotalPages:  totalPages,
		},
		List: data,
	}, nil
}

func (r *Repository) Detail(ctx context.Context, pg *sql.DB, employeeId uint) (*readmodel.EmployeeDetail, error) {
	data := &readmodel.EmployeeDetail{}

	var histories []byte

	if err := pg.QueryRowContext(ctx, employeeDetailQuery, employeeId).Scan(
		&data.Id,
		&data.Name,
		&data.Email,
		&data.Phone,
		&data.Gender,
		&data.Address,
		&data.JoinDate,
		&data.Legal.IdentityNumber,
		&data.Legal.TaxNumber,
		&data.Legal.HealthInsuranceNumber,
		&data.Legal.LaborInsuranceNumber,
		&histories,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrEmployeeNotFound
		}

		return nil, err
	}

	if err := json.Unmarshal(histories, &data.Histories); err != nil {
		return nil, err
	}

	return data, nil
}

func New() *Repository {
	return &Repository{}
}
