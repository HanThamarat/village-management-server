package request


type CreateTokenOPN struct {
	Name      			string 			`validate:"required" json:"name"`
	Number    			string 			`validate:"required" json:"number"`
	ExpirationMonth		uint 			`validate:"required" json:"expirationMonth"`
	ExpirationYear      uint   			`validate:"required" json:"expirationYear"`
	City				string   		`validate:"required" json:"city"`
	PostalCode			string   		`validate:"required" json:"postalCode"`
	SecurityCode		string   		`validate:"required" json:"securityCode"`
}