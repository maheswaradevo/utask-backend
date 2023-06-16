package service

import (
	"context"

	"github.com/maheswaradevo/utask-backend/internal/authentications"
	"github.com/maheswaradevo/utask-backend/internal/calendar"
	"github.com/maheswaradevo/utask-backend/internal/models"
	"github.com/maheswaradevo/utask-backend/pkg/common/constants"
	"github.com/maheswaradevo/utask-backend/pkg/common/helpers"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type calendarService struct {
	logger        *zap.Logger
	calendarRest  calendar.CalendarRestRepository
	authRedisRepo authentications.GoogleOauthRedisRepository
}

func NewCalendarService(calendarRest calendar.CalendarRestRepository, authRedisRepo authentications.GoogleOauthRedisRepository, logger *zap.Logger) *calendarService {
	return &calendarService{
		calendarRest:  calendarRest,
		authRedisRepo: authRedisRepo,
		logger:        logger,
	}
}

func (s *calendarService) GetEvent(ctx context.Context) (*models.EventList, error) {
	var redisKey = helpers.CacheWithPrefix(constants.CacheAuthGoogle, constants.CacheGoogleLogin)
	data, err := s.authRedisRepo.FindByKey(ctx, redisKey)
	if err != nil {
		s.logger.Error("error: %v", zap.Error(err))
		return nil, err
	}

	eventList, err := s.calendarRest.GetEventList(ctx, &oauth2.Token{
		AccessToken:  data.AccessToken,
		TokenType:    data.TokenType,
		RefreshToken: data.RefreshToken,
		Expiry:       data.Expiry,
	})
	if err != nil {
		s.logger.Error("error: %v", zap.Error(err))
		return nil, err
	}
	return eventList, nil
}

func (s *calendarService) GetEventByID(ctx context.Context, eventId string) (*models.EventResource, error) {
	var redisKey = helpers.CacheWithPrefix(constants.CacheAuthGoogle, constants.CacheGoogleLogin)
	data, err := s.authRedisRepo.FindByKey(ctx, redisKey)
	if err != nil {
		s.logger.Error("error: %v", zap.Error(err))
		return nil, err
	}

	event, err := s.calendarRest.GetEventByID(ctx, eventId, &oauth2.Token{
		AccessToken:  data.AccessToken,
		TokenType:    data.TokenType,
		RefreshToken: data.RefreshToken,
		Expiry:       data.Expiry,
	})
	if err != nil {
		s.logger.Error("error: %v", zap.Error(err))
		return nil, err
	}
	return event, nil
}
