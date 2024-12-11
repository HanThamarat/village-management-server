package request

type RegisterReq struct {
	Name     		string `validate:"required" json:"name"`
	Username     	string `validate:"required" json:"username"`
	Password     	string `validate:"required" json:"password"`
	RoleReq 		uint 	`validate:"required" json:"role"`
}

type SigninReq struct {
	Username         string `validate:"required" json:"username"`
	Password         string `validate:"required" json:"password"`
}