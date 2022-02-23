package service

import (
	"context"
)

// QueryEvents queries logs
func (s *service) QueryEvents(ctx context.Context, eventType string, query map[string]string) ([]Event, error) {
	return s.datastoreRepository.QueryEvents(ctx, eventType, query)
}
