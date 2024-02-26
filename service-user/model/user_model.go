package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       string `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null" json:"password"`
}

func (e *User) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New().String()
	return
}
