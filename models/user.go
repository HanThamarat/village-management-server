package models

import (
	"time"
)

type Users struct {
	ID                uint
	Name              string
	Username          string
	Password          string
	Roles			  uint
	Role              Roles `gorm:"foreignKey:Roles"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
