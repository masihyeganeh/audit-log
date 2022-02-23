package service

import (
	"context"
	"github.com/masihyeganeh/audit-log/internal/auth"
	"math"
	"runtime"
	"time"
)

// Service interface provides membership business logic functions
type Service interface {
	// StartBackgroundWorkers starts background workers
	StartBackgroundWorkers(ctx context.Context)
	// LogEvent logs the event
	LogEvent(ctx context.Context, event Event) error
	// QueryEvents queries logs
	QueryEvents(ctx context.Context, eventType string, query map[string]string) ([]Event, error)
	// Login user login
	Login(ctx context.Context, request LoginRequest) (string, error)
}

type service struct {
	datastoreRepository     DatastoreRepository
	authentication          auth.Auth
	channel                 chan Event
	workersCount            int
	eventsBufferSize        int
	maxEventsBufferDuration time.Duration
}

type Config struct {
	ChannelSize             int           `yaml:"SERVICE_CONFIG_CHANNEL_SIZE"`
	MaxWorkers              int           `yaml:"SERVICE_CONFIG_MAX_WORKERS"`
	EventsBufferSize        int           `yaml:"SERVICE_CONFIG_EVENTS_BUFFER_SIZE"`
	MaxEventsBufferDuration time.Duration `yaml:"SERVICE_CONFIG_MAX_EVENTS_BUFFER_DURATION"`
}

// CreateService creates an instance of membership service with the necessary dependencies
func CreateService(datastoreRepository DatastoreRepository, authentication auth.Auth, config Config) Service {
	return &service{
		datastoreRepository:     datastoreRepository,
		authentication:          authentication,
		channel:                 make(chan Event, config.ChannelSize),
		workersCount:            int(math.Min(float64(runtime.NumCPU()-1), float64(config.MaxWorkers))),
		eventsBufferSize:        config.EventsBufferSize,
		maxEventsBufferDuration: config.MaxEventsBufferDuration,
	}
}
