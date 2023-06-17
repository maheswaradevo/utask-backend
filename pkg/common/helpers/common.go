package helpers

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/maheswaradevo/utask-backend/internal/models"
	"golang.org/x/oauth2"
)

func HandleLogin(w http.ResponseWriter, r *http.Request, oauthConf *oauth2.Config, oauthStateString string) {
	URL, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Error("Parse: " + err.Error())
	}
	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	parameters.Add("access_type", "offline")
	parameters.Add("prompt", "consent")
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallbackFromGoogle(w http.ResponseWriter, r *http.Request, oauthConfGl *oauth2.Config, oauthStateStringGl string) *oauth2.Token {
	state := r.FormValue("state")

	if state != oauthStateStringGl {
		log.Info("invalid oauth state, expected " + oauthStateStringGl + ", got " + state + "\n")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil
	}

	code := r.FormValue("code")

	if code == "" {
		log.Warn("Code not found..")
		w.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			w.Write([]byte("User has denied Permission.."))
		}
	} else {
		token, err := oauthConfGl.Exchange(context.Background(), code)
		if err != nil {
			log.Error("oauthConfGl.Exchange() failed with " + err.Error() + "\n")
			return nil
		}
		var authData models.Authentication
		authData.AccessToken = token.AccessToken
		authData.RefreshToken = token.RefreshToken
		authData.Expiry = token.Expiry
		authData.TokenType = token.TokenType
		return token

	}
	return nil
}
