package service

import "time"

// Event Domain Model for events
type Event struct {
	EventTime    time.Time
	EventType    string
	CommonField1 string
	CommonField2 string
	Fields       map[string]string
}

// LoginRequest login request
type LoginRequest struct {
	Username string
	Password string
}
