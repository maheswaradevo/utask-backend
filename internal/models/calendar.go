package models

import "time"

type EventList struct {
	Kind             string    `json:"kind"`
	Etag             string    `json:"etag"`
	Summary          string    `json:"summary"`
	Description      string    `json:"description"`
	Updated          time.Time `json:"updated"`
	TimeZone         string    `json:"timeZone"`
	AccessRole       string    `json:"accessRole"`
	DefaultReminders []struct {
		Method  string `json:"method"`
		Minutes int    `json:"minutes"`
	} `json:"defaultReminders"`
	NextPageToken string          `json:"nextPageToken"`
	NextSyncToken string          `json:"nextSyncToken"`
	Items         []EventResource `json:"items"`
}

type EventResource struct {
	Kind     string    `json:"kind"`
	Etag     string    `json:"etag"`
	ID       string    `json:"id"`
	Status   string    `json:"status"`
	HTMLLink string    `json:"htmlLink"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Summary  string    `json:"summary"`
	Creator  struct {
		Email       string `json:"email"`
		DisplayName string `json:"displayName"`
		Self        bool   `json:"self"`
	} `json:"creator"`
	Organizer struct {
		Email       string `json:"email"`
		DisplayName string `json:"displayName"`
		Self        bool   `json:"self"`
	} `json:"organizer"`
	Start struct {
		DateTime time.Time `json:"dateTime"`
		TimeZone string    `json:"timeZone"`
	} `json:"start"`
	End struct {
		DateTime time.Time `json:"dateTime"`
		TimeZone string    `json:"timeZone"`
	} `json:"end"`
	ICalUID   string `json:"iCalUID"`
	Sequence  int    `json:"sequence"`
	Reminders struct {
		UseDefault bool `json:"useDefault"`
	} `json:"reminders"`
	EventType string `json:"eventType"`
}
