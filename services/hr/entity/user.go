package entity

import (
	"enterprise_resource_planning/internal/shared/errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string    `json:"username,omitempty" gorm:"uniqueIndex:idx_active_username,where:deleted_at is null;type:varchar(100)"`
	Password   string    `json:"password,omitempty" gorm:"type:text;not null"`
	Role       string    `json:"role" gorm:"type:varchar(25);not null;check:,role <> 'superadmin','manager','staff'"`
	EmployeeID uint      `json:"employee_id,omitempty" gorm:"uniqueIndex:idx_active_user_employee,where:deleted_at is null"`
	Employee   *Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
}

func (u *User) GenerateHashPassowrd() error {
	if u.Password == "" {
		return errors.ErrBlankPassword
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MaxCost)

	if err != nil {
		return err
	}

	u.Password = string(hashPassword)

	return nil
}

func (u *User) CompareHashAndPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return errors.ErrInvalidPassword
	}

	return nil
}
