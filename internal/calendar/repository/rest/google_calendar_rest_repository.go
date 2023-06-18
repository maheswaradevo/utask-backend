package rest

import (
	"context"
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
	tMax := time.Now().Add(45 * time.Hour).Format(time.RFC3339)

	tMin := time.Now().Add(-720 * time.Hour).Format(time.RFC3339)

	svc, err := calendar.NewService(context.Background(), option.WithHTTPClient(config.CalendarServiceClient(token)))
	if err != nil {
		c.logger.Error("error: %v", zap.Error(err))
		return nil, err
	}
	events, err := svc.Events.List("primary").TimeMin(tMin).TimeMax(tMax).MaxResults(100).Do()
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
	if len(events.Items) == 0 {
		c.logger.Info("no items to show")
	} else {
		for item, eventResource := range events.Items {
			overrides := make([]models.Overrides, len(events.Items[item].Reminders.Overrides))

			eventResources[item].Kind = eventResource.Kind
			eventResources[item].Etag = eventResource.Etag
			eventResources[item].ID = eventResource.Id
			eventResources[item].Status = eventResource.Status
			eventResources[item].HTMLLink = eventResource.HtmlLink
			eventResources[item].Created = eventResource.Created
			eventResources[item].Summary = eventResource.Summary

			eventResources[item].Start.DateTime = eventResource.Start.DateTime
			eventResources[item].Start.TimeZone = eventResource.Start.TimeZone

			eventResources[item].End.DateTime = eventResource.End.DateTime
			eventResources[item].End.TimeZone = eventResource.End.TimeZone

			eventResources[item].Creator.DisplayName = eventResource.Creator.DisplayName
			eventResources[item].Creator.Email = eventResource.Creator.Email

			eventResources[item].Organizer.DisplayName = eventResource.Organizer.DisplayName
			eventResources[item].Organizer.Email = eventResource.Organizer.Email

			eventResources[item].Reminders.UseDefault = eventResource.Reminders.UseDefault

			for i, override := range eventResource.Reminders.Overrides {
				overrides[i].Method = override.Method
				overrides[i].Minutes = int(override.Minutes)
			}
			eventResources[item].Reminders.Overrides = overrides
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

	overrides := make([]models.Overrides, 2)
	for item, override := range event.Reminders.Overrides {
		overrides[item].Method = override.Method
		overrides[item].Minutes = int(override.Minutes)
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
		Etag:        event.Etag,
		HTMLLink:    event.HtmlLink,
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
			Overrides:  overrides,
		},
		Organizer: models.Organizer{
			Email:       event.Organizer.Email,
			DisplayName: event.Organizer.DisplayName,
		},
	}

	return &eventDetail, err
}
