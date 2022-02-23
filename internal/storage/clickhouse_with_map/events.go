package clickhouse

import (
	"github.com/masihyeganeh/audit-log/internal/service"
	"time"
)

type event struct {
	EventTime    time.Time         `json:"event_time"`
	EventType    string            `json:"event_type"`
	CommonField1 string            `json:"common_field_1"`
	CommonField2 string            `json:"common_field_2"`
	Fields       map[string]string `json:"fields"`
}

func (e event) toServiceModel() service.Event {
	return service.Event{
		EventTime:    e.EventTime,
		EventType:    e.EventType,
		CommonField1: e.CommonField1,
		CommonField2: e.CommonField2,
		Fields:       e.Fields,
	}
}

type logEvents struct {
	Events []event
}

func (l logEvents) toServiceModel() []service.Event {
	events := make([]service.Event, len(l.Events))

	for i, event := range l.Events {
		events[i] = event.toServiceModel()
	}

	return events
}
