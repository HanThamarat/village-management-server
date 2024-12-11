package request

type RoleReq struct {
	Name_TH string `json:"Name_TH"`
	Name_EN string `json:"Name_EN"`
	Active  bool   `json:"Active"`
}