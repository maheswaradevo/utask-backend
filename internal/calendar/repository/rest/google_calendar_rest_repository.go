package rest

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/maheswaradevo/utask-backend/pkg/config"
	"golang.org/x/oauth2"
	calendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type calendarRestRepository struct {
}

func NewCalendarRestRepository() calendarRestRepository {
	return calendarRestRepository{}
}

func (c *calendarRestRepository) GetEventList(ctx context.Context, token *oauth2.Token) {
	t := time.Now().Format(time.RFC3339)

	svc, err := calendar.NewService(context.Background(), option.WithHTTPClient(config.CalendarServiceClient(token)))
	if err != nil {
		log.Printf("error ni: %v", err)
	}
	events, err := svc.Events.List("primary").TimeMin(t).MaxResults(10).Do()

	if err != nil {
		log.Printf("ada error: %v", err)
		return
	}

	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%v (%v)\n", item.Summary, date)
		}
	}
}
