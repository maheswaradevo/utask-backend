package service

import (
	"context"
	"net/http"
	"time"

	"github.com/maheswaradevo/utask-backend/internal/authentications"
	"github.com/maheswaradevo/utask-backend/internal/models"
	"github.com/maheswaradevo/utask-backend/pkg/common/constants"
	"github.com/maheswaradevo/utask-backend/pkg/common/helpers"
	"github.com/maheswaradevo/utask-backend/pkg/config"

	"golang.org/x/oauth2"
)

type googleOauthService struct {
	googleRedisRepository authentications.GoogleOauthRedisRepository
}

func NewGoogleOauthService(googleRedisRepository authentications.GoogleOauthRedisRepository) authentications.GoogleOauthService {
	return &googleOauthService{googleRedisRepository: googleRedisRepository}
}

func (s *googleOauthService) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	helpers.HandleLogin(w, r, config.OauthConfGl, config.OauthStateStringGl)
}

func (s *googleOauthService) CallBackFromGoogle(w http.ResponseWriter, r *http.Request) *oauth2.Token {
	token := helpers.CallbackFromGoogle(w, r, config.OauthConfGl, config.OauthStateStringGl)
	var redisKey = helpers.CacheWithPrefix(constants.CacheAuthGoogle, constants.CacheGoogleLogin)

	var expTime float64

	_, err := s.googleRedisRepository.FindByKey(context.Background(), redisKey)
	if err != nil && expTime == 0 {
		data := models.OauthRedisData{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			Key:          redisKey,
			ExpiredTime:  int(1 * time.Hour),
			Expiry:       token.Expiry,
			TokenType:    token.TokenType,
		}

		_, errSave := s.googleRedisRepository.Save(context.Background(), redisKey, &data)
		if errSave != nil {
			return nil
		}
	}
	return token
}
