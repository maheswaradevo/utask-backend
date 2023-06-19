package authentications

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

type GoogleOauthService interface {
	HandleGoogleLogin(w http.ResponseWriter, r *http.Request)
	CallBackFromGoogle(w http.ResponseWriter, r *http.Request, state, code string) (*oauth2.Token, error)
	Logout(ctx context.Context) (bool, error)
}
