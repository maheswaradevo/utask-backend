package notification

import "time"

type SendSMS struct {
	Destination string    `json:"destination"`
	Country     string    `json:"country"`
	Source      string    `json:"source"`
	Text        string    `json:"text"`
	Encoding    string    `json:"encoding"`
	Scheduled   time.Time `json:"scheduled"`
}

type ResponseError struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	ErrorID   string    `json:"errorId"`
	Timestamp time.Time `json:"timestamp"`
}
