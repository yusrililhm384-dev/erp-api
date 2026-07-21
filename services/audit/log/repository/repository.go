package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
}

func (r *Repository) SetUserId(tx *gorm.DB, userId int64) error {
	return tx.Exec("SELECT set_config('app.user_id', ?, true)", userId).Error
}

func New() *Repository {
	return &Repository{}
}
