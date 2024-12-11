package models

import (
    "time"
)

type VallageHouses struct {
    ID  uint
    HouseName *string
	HousePrice *float64
	WaterBill *float64
	ElectricityBill *float64
	OtherBill *float64
    Active *bool
	VallageOwnerShipID uint
	VallageOwnerShip VallageOwnerShips `gorm:"foreignKey:VallageOwnerShipID"`
	CreateBy *uint
	Create Users `gorm:"foreignKey:CreateBy"`
    CreatedAt time.Time
	UpdateBy *uint
	Update Users `gorm:"foreignKey:UpdateBy"`

	UpdatedAt time.Time
}