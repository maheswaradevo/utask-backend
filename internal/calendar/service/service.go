package service

import (
	"context"

	"github.com/maheswaradevo/utask-backend/internal/authentications"
	"github.com/maheswaradevo/utask-backend/internal/calendar"
	"github.com/maheswaradevo/utask-backend/pkg/common/constants"
	"github.com/maheswaradevo/utask-backend/pkg/common/helpers"
	"golang.org/x/oauth2"
)

type calendarService struct {
	calendarRest  calendar.CalendarRestRepository
	authRedisRepo authentications.GoogleOauthRedisRepository
}

func NewCalendarService(calendarRest calendar.CalendarRestRepository, authRedisRepo authentications.GoogleOauthRedisRepository) *calendarService {
	return &calendarService{
		calendarRest:  calendarRest,
		authRedisRepo: authRedisRepo,
	}
}

func (s *calendarService) GetEvent(ctx context.Context) {
	var redisKey = helpers.CacheWithPrefix(constants.CacheAuthGoogle, constants.CacheGoogleLogin)
	data, err := s.authRedisRepo.FindByKey(ctx, redisKey)
	if err != nil {
		return
	}

	s.calendarRest.GetEventList(ctx, &oauth2.Token{
		AccessToken:  data.AccessToken,
		TokenType:    data.TokenType,
		RefreshToken: data.RefreshToken,
		Expiry:       data.Expiry,
	})
}
