package ipcarta

type Network struct {
	Network          string `json:"network"`
	IsAnonymousProxy bool   `json:"isAnonymousProxy"`
	ContinentCode    string `json:"continentCode"`
	ContinentName    string `json:"continentName"`
	CountryISOCode   string `json:"countryISOCode"`
	CountryName      string `json:"countryName"`
}
