package request

type HouseReq struct {
	HouseName string `json:"house_name"`
	HousePrice float64 `json:"house_price"`
	WaterBill *float64	`json:"water_bill"`
	ElectricityBill *float64    `json:"electricity_bill"`
	OtherBill *float64    `json:"other_bill"`
	Active bool `json:"active"`
	VallageOwnerShipID uint `json:"vallage_owner_ship_id"`
}