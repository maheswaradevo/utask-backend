package models

type OauthRedisData struct {
	AccessToken  string
	RefreshToken string
	ExpiredTime  int
	Key          string
}
