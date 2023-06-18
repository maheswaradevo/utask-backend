package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/maheswaradevo/utask-backend/internal/authentications"
	"github.com/maheswaradevo/utask-backend/internal/calendar"
	"github.com/maheswaradevo/utask-backend/internal/models"
	"github.com/maheswaradevo/utask-backend/internal/notification"
	"github.com/maheswaradevo/utask-backend/pkg/common/constants"
	"github.com/maheswaradevo/utask-backend/pkg/common/helpers"
	"github.com/maheswaradevo/utask-backend/pkg/common/helpers/json"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type calendarService struct {
	logger         *zap.Logger
	calendarRest   calendar.CalendarRestRepository
	authRedisRepo  authentications.GoogleOauthRedisRepository
	wavecellClient notification.Client
}

func NewCalendarService(calendarRest calendar.CalendarRestRepository, authRedisRepo authentications.GoogleOauthRedisRepository, wavecellClient notification.Client, logger *zap.Logger) *calendarService {
	return &calendarService{
		calendarRest:   calendarRest,
		authRedisRepo:  authRedisRepo,
		wavecellClient: wavecellClient,
		logger:         logger,
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

func (s *calendarService) SendSMS(ctx context.Context, eventId, phoneNumber string) (result models.Notification, err error) {
	var redisKey = helpers.CacheWithPrefix(constants.CacheAuthGoogle, constants.CacheGoogleLogin)
	data, err := s.authRedisRepo.FindByKey(ctx, redisKey)
	if err != nil {
		s.logger.Error("error: %v", zap.Error(err))
		return
	}
	event, err := s.calendarRest.GetEventByID(ctx, eventId, &oauth2.Token{
		AccessToken:  data.AccessToken,
		TokenType:    data.TokenType,
		RefreshToken: data.RefreshToken,
		Expiry:       data.Expiry,
	})
	if err != nil {
		s.logger.Error("error: %v", zap.Error(err))
		return
	}
	var reminder time.Duration

	if event.Reminders.UseDefault {
		reminder = 30 * time.Minute
	} else {
		reminder = time.Duration(event.Reminders.Overrides[0].Minutes * int(time.Minute))
	}
	startDate := helpers.ParseDate(&event.Start.DateTime, helpers.DateLayoutTimeZone)
	notify := startDate.Add(-reminder * time.Minute)

	response, errSendSMS := s.wavecellClient.SendSMS(helpers.Env("SMS_API_URL"), notification.SendSMS{
		Destination: phoneNumber,
		Country:     "ID",
		Source:      "u-Task",
		Text:        fmt.Sprintf("Ingat ada kegiatan %v", event.Summary),
		Encoding:    "AUTO",
		Scheduled:   notify,
	})
	if errSendSMS != nil {
		err = errSendSMS
		s.logger.Sugar().Errorf("error: %v", errSendSMS)
		return
	}
	if response.StatusCode() != http.StatusOK {
		var respError notification.ResponseError
		if errDecodeJson := json.FromJSON(response.Body(), &respError); errDecodeJson != nil {
			err = errDecodeJson
			return
		}
		err = errors.New(respError.Message)
		return
	}
	result.EventId = event.ID
	result.PhoneNumber = phoneNumber
	return
}
