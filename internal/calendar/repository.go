package calendar

import (
	"context"

	"github.com/maheswaradevo/utask-backend/internal/models"
	"golang.org/x/oauth2"
)

type CalendarRestRepository interface {
	GetEventList(ctx context.Context, token *oauth2.Token) (*models.EventList, error)
	// GetEventByID(ctx context.Context, eventId string, token *oauth2.Token) (*models.EventResource, error)
}

type CalendarRepository interface{}
