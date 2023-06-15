package calendar

import (
	"context"

	"golang.org/x/oauth2"
)

type CalendarRestRepository interface {
	GetEventList(ctx context.Context, token *oauth2.Token)
}

type CalendarRepository interface{}
