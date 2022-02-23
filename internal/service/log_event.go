package service

import (
	"context"
	"time"
)

// LogEvent logs the event
func (s *service) LogEvent(ctx context.Context, event Event) error {
	event.EventTime = time.Now()
	s.channel <- event
	return nil
}
