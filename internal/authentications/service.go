package authentications

import (
	"net/http"

	"golang.org/x/oauth2"
)

type GoogleOauthService interface {
	HandleGoogleLogin(w http.ResponseWriter, r *http.Request)
	CallBackFromGoogle(w http.ResponseWriter, r *http.Request) *oauth2.Token
}
