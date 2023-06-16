package rest

import (
	"context"
	"fmt"
	"time"

	"github.com/maheswaradevo/utask-backend/internal/models"
	"github.com/maheswaradevo/utask-backend/pkg/config"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	calendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type calendarRestRepository struct {
	logger *zap.Logger
}

func NewCalendarRestRepository(logger *zap.Logger) calendarRestRepository {
	return calendarRestRepository{logger: logger}
}

func (c *calendarRestRepository) GetEventList(ctx context.Context, token *oauth2.Token) (*models.EventList, error) {
	t := time.Now().Format(time.RFC3339)

	svc, err := calendar.NewService(context.Background(), option.WithHTTPClient(config.CalendarServiceClient(token)))
	if err != nil {
		c.logger.Error("error: %v", zap.Error(err))
		return nil, err
	}
	events, err := svc.Events.List("primary").TimeMin(t).MaxResults(10).Do()

	if err != nil {
		c.logger.Error("error: %v", zap.Error(err))
		return nil, err
	}

	defaultReminders := make([]models.DefaultReminder, len(events.DefaultReminders))
	for item, defaultReminder := range events.DefaultReminders {
		defaultReminders[item].Method = defaultReminder.Method
		defaultReminders[item].Minutes = int(defaultReminder.Minutes)
	}

	eventResources := make([]models.EventResource, len(events.Items))
	fmt.Printf("len(events.Items): %v\n", len(events.Items))
	if len(events.Items) == 0 {
		c.logger.Info("no items to show")
	} else {
		for item, eventResource := range events.Items {
			eventResources[item].Kind = eventResource.Kind
			eventResources[item].Etag = eventResource.Etag
			eventResources[item].ID = eventResource.Id
			eventResources[item].Status = eventResource.Status
			eventResources[item].HTMLLink = eventResource.HtmlLink
			// createdTime, err := time.Parse("2006-01-02T15:04:05Z", eventResource.Created)
			// if err != nil {
			// 	c.logger.Error("error: %v", zap.Error(err))
			// 	return nil, err
			// }
			eventResources[item].Created = eventResource.Created
			eventResources[item].Summary = eventResource.Summary
			// startDateTime, err := time.Parse("2006-01-02", eventResource.Start.DateTime)
			// if err != nil {
			// 	c.logger.Error("error: %v", zap.Error(err))
			// 	return nil, err
			// }
			eventResources[item].Start.DateTime = eventResource.Start.DateTime
			eventResources[item].Start.TimeZone = eventResource.Start.TimeZone
			// endDateTime, err := time.Parse("2006-01-02T15:04:05Z", eventResource.End.DateTime)
			// if err != nil {
			// 	c.logger.Error("error: %v", zap.Error(err))
			// 	return nil, err
			// }
			eventResources[item].End.DateTime = eventResource.End.DateTime
			eventResources[item].End.TimeZone = eventResource.End.TimeZone

			eventResources[item].Creator.DisplayName = eventResource.Creator.DisplayName
			eventResources[item].Creator.Email = eventResource.Creator.Email

			eventResources[item].Organizer.DisplayName = eventResource.Organizer.DisplayName
			eventResources[item].Organizer.Email = eventResource.Organizer.Email
		}
	}

	var eventList = models.EventList{
		Kind:             events.Kind,
		Summary:          events.Summary,
		Description:      events.Summary,
		Updated:          events.Updated,
		TimeZone:         events.TimeZone,
		NextPageToken:    events.NextPageToken,
		DefaultReminders: defaultReminders,
		NextSyncToken:    events.NextSyncToken,
		Items:            eventResources,
	}

	return &eventList, nil
}

func (c *calendarRestRepository) GetEventByID(ctx context.Context, eventId string, token *oauth2.Token) (*models.EventResource, error) {
	svc, err := calendar.NewService(context.Background(), option.WithHTTPClient(config.CalendarServiceClient(token)))
	if err != nil {
		c.logger.Error("error: %v", zap.Error(err))
		return nil, err
	}
	event, err := svc.Events.Get("primary", eventId).Do()
	if err != nil {
		c.logger.Error("error: %v", zap.Error(err))
		return nil, err
	}

	start := models.StartTime{
		DateTime: event.Start.DateTime,
		TimeZone: event.Start.TimeZone,
	}

	end := models.EndTime{
		DateTime: event.End.DateTime,
		TimeZone: event.End.TimeZone,
	}
	eventDetail := models.EventResource{
		Kind:        event.Kind,
		Summary:     event.Summary,
		HangoutLink: event.HangoutLink,
		ID:          event.Id,
		Created:     event.Created,
		Updated:     event.Updated,
		Start:       start,
		End:         end,
		Creator: models.Creator{
			Email:       event.Creator.Email,
			DisplayName: event.Creator.DisplayName,
		},
		Reminders: models.Reminders{
			UseDefault: event.Reminders.UseDefault,
		},
	}

	return &eventDetail, err
}
