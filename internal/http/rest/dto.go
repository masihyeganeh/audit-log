package rest

import "time"

type LogRequest struct {
	EventType    string            `json:"event_type"`
	CommonField1 string            `json:"common_field_1"`
	CommonField2 string            `json:"common_field_2"`
	Fields       map[string]string `json:"fields"`
}

type LogResponse struct {
	EventTime    time.Time         `json:"event_time"`
	EventType    string            `json:"event_type"`
	CommonField1 string            `json:"common_field_1"`
	CommonField2 string            `json:"common_field_2"`
	Fields       map[string]string `json:"fields"`
}

type QueryRequest struct {
	EventType string            `json:"event_type"`
	Filters   map[string]string `json:"filters"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
