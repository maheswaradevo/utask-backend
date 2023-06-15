package service

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/maheswaradevo/utask-backend/internal/authentications"
	"github.com/maheswaradevo/utask-backend/internal/models"
	"github.com/maheswaradevo/utask-backend/pkg/common/constants"
	"github.com/maheswaradevo/utask-backend/pkg/common/helpers"
	"github.com/maheswaradevo/utask-backend/pkg/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfGl = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:3000/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateStringGl = ""
)

type googleOauthService struct {
	googleRedisRepository authentications.GoogleOauthRedisRepository
}

func NewGoogleOauthService(googleRedisRepository authentications.GoogleOauthRedisRepository) authentications.GoogleOauthService {
	return &googleOauthService{googleRedisRepository: googleRedisRepository}
}

func (s *googleOauthService) InitializeOauthGoogle() {
	cfg := config.GetConfig()
	oauthConfGl.ClientID = cfg.GoogleClientID
	oauthConfGl.ClientSecret = cfg.GoogleClientSecret
	oauthStateStringGl = cfg.GoogleStateString
}

func (s *googleOauthService) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	helpers.HandleLogin(w, r, oauthConfGl, oauthStateStringGl)
}

func (s *googleOauthService) CallBackFromGoogle(w http.ResponseWriter, r *http.Request) {
	log.Info("Callback-gl..")

	state := r.FormValue("state")
	log.Info(state)
	if state != oauthStateStringGl {
		log.Info("invalid oauth state, expected " + oauthStateStringGl + ", got " + state + "\n")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	log.Info(code)

	if code == "" {
		log.Warn("Code not found..")
		w.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			w.Write([]byte("User has denied Permission.."))
		}
		// User has denied access..
		// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		token, err := oauthConfGl.Exchange(context.Background(), code)
		if err != nil {
			log.Error("oauthConfGl.Exchange() failed with " + err.Error() + "\n")
			return
		}
		log.Info("TOKEN>> AccessToken>> " + token.AccessToken)
		log.Info("TOKEN>> Expiration Time>> ", token.Expiry.Hour())
		log.Info("TOKEN>> RefreshToken>> " + token.RefreshToken)

		var redisKey = helpers.CacheWithPrefix(constants.CacheAuthGoogle, constants.CacheGoogleLogin)

		var expTime float64

		_, err = s.googleRedisRepository.FindByKey(context.Background(), redisKey)
		if err != nil && expTime == 0 {
			data := models.OauthRedisData{
				AccessToken:  token.AccessToken,
				RefreshToken: token.RefreshToken,
				Key:          redisKey,
				ExpiredTime:  int(1 * time.Hour),
			}

			_, errSave := s.googleRedisRepository.Save(context.Background(), redisKey, &data)
			if errSave != nil {
				return
			}
		}
		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
		if err != nil {
			log.Error("Get: " + err.Error() + "\n")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		defer resp.Body.Close()

		response, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error("ReadAll: " + err.Error() + "\n")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		log.Info("parseResponseBody: " + string(response) + "\n")

		w.Write([]byte("Hello, I'm protected\n"))
		w.Write([]byte(string(response)))
		return
	}
}
