package linkdtos

type TokenLink struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	Audience    string `json:"audience"`
	Expires_in  string `json:"expires_in"`
}
