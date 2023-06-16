package models

type EventList struct {
	Kind             string            `json:"kind"`
	Etag             string            `json:"etag"`
	Summary          string            `json:"summary"`
	Description      string            `json:"description"`
	Updated          string            `json:"updated"`
	TimeZone         string            `json:"timeZone"`
	AccessRole       string            `json:"accessRole"`
	DefaultReminders []DefaultReminder `json:"defaultReminders"`
	NextPageToken    string            `json:"nextPageToken"`
	NextSyncToken    string            `json:"nextSyncToken"`
	Items            []EventResource   `json:"items"`
}

type EventResource struct {
	Kind      string    `json:"kind"`
	Etag      string    `json:"etag"`
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	HTMLLink  string    `json:"htmlLink"`
	Created   string    `json:"created"`
	Updated   string    `json:"updated"`
	Summary   string    `json:"summary"`
	Creator   Creator   `json:"creator"`
	Organizer Organizer `json:"organizer"`
	Start     StartTime `json:"start"`
	End       EndTime   `json:"end"`
	ICalUID   string    `json:"iCalUID"`
	Sequence  int       `json:"sequence"`
	Reminders Reminders `json:"reminders"`
	EventType string    `json:"eventType"`
}

type DefaultReminder struct {
	Method  string `json:"method"`
	Minutes int    `json:"minutes"`
}

type StartTime struct {
	DateTime string `json:"dateTime"`
	TimeZone string `json:"timeZone"`
}

type EndTime struct {
	DateTime string `json:"dateTime"`
	TimeZone string `json:"timeZone"`
}

type Creator struct {
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	Self        bool   `json:"self"`
}

type Organizer struct {
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	Self        bool   `json:"self"`
}

type Reminders struct {
	UseDefault bool `json:"useDefault"`
}
