package request

type GenAccessTokenHeader struct {
	ResourceOwnerId	string `json:"resourceOwnerId"`
	RequestUId	string `json:"requestUId"`
	Acceptlanguage	string `json:"accept-language"`
}

type GenAccessTokenBody struct {
	ApplicationKey	string `json:"applicationKey"`
	ApplicationSecret string `json:"applicationSecret"`
}