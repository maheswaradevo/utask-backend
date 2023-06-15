package models

import "time"

type OauthRedisData struct {
	AccessToken  string
	RefreshToken string
	ExpiredTime  int
	Key          string
	TokenType    string
	Expiry       time.Time
}

type Authentication struct {
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
	TokenType    string
}
