package service

import "context"

// DatastoreRepository provides access to datastore repository
type DatastoreRepository interface {
	// Close closes connection to datastore
	Close() error
	// LogEvents writes log to datastore
	LogEvents(ctx context.Context, event []Event) error
	// QueryEvents queries logs from datastore
	QueryEvents(ctx context.Context, eventType string, filters map[string]string) ([]Event, error)
	// AddUser creates a user with specific access
	AddUser(ctx context.Context, username, hashedPassword, salt string, hasReadAccess, hasWriteAccess bool) error
	// FindUser returns user data with username
	FindUser(ctx context.Context, username string) (string, string, bool, bool, error)
}
