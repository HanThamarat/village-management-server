package models

import (
	"time"
)

type VallageOwnerShips struct {
	ID            uint `gorm:"primaryKey"`
	VillageName   *string
	LegalEntityID uint
	LegalEntity   Users `gorm:"foreignKey:LegalEntityID"`
	Active        *string
	CreateByID    uint
	Create        Users `gorm:"foreignKey:CreateByID"`
	UpdateByID    uint
	Update        Users `gorm:"foreignKey:UpdateByID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
