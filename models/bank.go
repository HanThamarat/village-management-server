package models

import (
	"time"
)

type BankCredentials struct {
	ID uint
	BankName *string
	AppName *string
	API_KEY *string
	API_SECRET *string
	Biller_ID *uint
	Merchant_ID *uint
	Terminal_ID *uint
	Active 	*bool
	CreatedAt time.Time
	UpdatedAt time.Time
}