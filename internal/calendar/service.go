package calendar

import (
	"context"
)

type CalendarService interface {
	GetEvent(ctx context.Context)
}
