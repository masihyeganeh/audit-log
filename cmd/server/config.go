package main

import "github.com/masihyeganeh/audit-log/internal/service"

// Config holds the app master configuration
type Config struct {
	ServerAddress               string         `yaml:"SERVER_ADDR"`
	DataStorage                 string         `yaml:"DATA_STORAGE"`
	DataStorageConnectionString string         `yaml:"CONNECTION_STRING"`
	ServiceConfig               service.Config `yaml:"SERVICE_CONFIG"`
	JwtSecret                   string         `yaml:"JWT_SECRET"`
}
