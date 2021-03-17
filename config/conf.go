package config

const (
	AppId       = "101941144"
	AppKey      = "006ba718548af07ea08fe595aab0ccae"
	RedirectURI = "http://qqt.feibooks.com/qq/callback"
)

type PrivateInfo struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
}