package calendar

import (
	"context"

	"github.com/maheswaradevo/utask-backend/internal/models"
)

type CalendarService interface {
	GetEvent(ctx context.Context) (*models.EventList, error)
	GetEventByID(ctx context.Context, eventId string) (*models.EventResource, error)
}
