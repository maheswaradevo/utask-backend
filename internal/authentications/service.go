package authentications

import "net/http"

type GoogleOauthService interface {
	InitializeOauthGoogle()
	HandleGoogleLogin(w http.ResponseWriter, r *http.Request)
}
