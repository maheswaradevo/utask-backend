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

func (s *googleOauthService) CallBackFromGoogle(w http.ResponseWriter, r *http.Request, state, code string) (*oauth2.Token, error) {
	token, err := helpers.CallbackFromGoogle(w, r, config.OauthConfGl, state, code, config.OauthStateStringGl)
	if err != nil {
		return nil, err
	}
	var redisKey = helpers.CacheWithPrefix(constants.CacheAuthGoogle, constants.CacheGoogleLogin)

	var expTime float64

	_, err = s.googleRedisRepository.FindByKey(context.Background(), redisKey)
	if err != nil && expTime == 0 {
		data := models.OauthRedisData{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			Key:          redisKey,
			ExpiredTime:  int(8 * time.Hour),
			Expiry:       token.Expiry,
			TokenType:    token.TokenType,
		}

		_, errSave := s.googleRedisRepository.Save(context.Background(), redisKey, &data)
		if errSave != nil {
			return nil, errSave
		}
	}
	return token, nil
}

func (s *googleOauthService) Logout(ctx context.Context) (bool, error) {
	var redisKey = helpers.CacheWithPrefix(constants.CacheAuthGoogle, constants.CacheGoogleLogin)
	data, err := s.googleRedisRepository.FindByKey(ctx, redisKey)
	if err != nil {
		return false, err
	}
	url := "https://accounts.google.com/o/oauth2/revoke?token="

	_, err = http.Get(url + data.AccessToken)
	if err != nil {
		return false, err
	}
	_ = s.googleRedisRepository.DeleteKey(ctx, redisKey)
	return true, nil
}
