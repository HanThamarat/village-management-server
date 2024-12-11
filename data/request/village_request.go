package request

type VillageOwerReq struct {
	VillageName string `validate:"required" json:"villageName"`
	// LegalEntity uint `validate:"required" json:"legalEntity"`
}