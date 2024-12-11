package models

import (
	"time"
)

type Roles struct {
	ID uint
	Name_TH *string
	Name_EN *string
	Active 	*bool
	CreatedAt time.Time
	UpdatedAt time.Time
}