package models

type Notification struct {
	PhoneNumber string `json:"phone_number"`
	EventId     string `json:"event_id"`
}

type NotificationRequest struct {
	PhoneNumber string `json:"phone_number"`
}
