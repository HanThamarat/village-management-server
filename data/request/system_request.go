package request

type BankReq struct {
	BankName string `validate:"required" json:"bankName"`
	AppName	string `validate:"required" json:"appName"`
	ApiKey	string `validate:"required" json:"apiKey"`
	ApiSecret string `validate:"required" json:"apiSecret"`
	BillerID uint `validate:"required" json:"billerID"`
	MerchantID uint `validate:"required" json:"merchantID"`
	TerminalID uint `validate:"required" json:"terminalID"`
}