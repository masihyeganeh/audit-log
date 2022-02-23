package main

import (
	"context"
	"github.com/masihyeganeh/audit-log/internal/auth"
	"github.com/masihyeganeh/audit-log/internal/http/rest"
	"github.com/masihyeganeh/audit-log/internal/service"
	clickhouseWithMap "github.com/masihyeganeh/audit-log/internal/storage/clickhouse_with_map"
	clickhouseWithNested "github.com/masihyeganeh/audit-log/internal/storage/clickhouse_with_nested"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type server struct {
	sync.WaitGroup
	config         *Config
	httpServer     *http.Server
	restHandler    *rest.Handler
	datastore      service.DatastoreRepository
	authentication auth.Auth
	service        service.Service
}

// NewServer Creates a new instance of server application
func NewServer(cfg *Config) *server {
	return &server{
		config: cfg,
	}
}

// Initialize is responsible for app initialization and wrapping required dependencies
func (s *server) Initialize(ctx context.Context) error {
	var datastoreRepository service.DatastoreRepository
	var isFirstRun bool
	var err error

	// Initialize Auth Service
	authentication := auth.New(s.config.JwtSecret)

	// Initialize Datastore Repository
	switch s.config.DataStorage {
	case "clickhouse_with_map":
		datastoreRepository, isFirstRun, err = clickhouseWithMap.ConnectAndCreateRepository(ctx, s.config.DataStorageConnectionString)
	case "clickhouse_with_nested":
		datastoreRepository, isFirstRun, err = clickhouseWithNested.ConnectAndCreateRepository(ctx, s.config.DataStorageConnectionString)
	default:
		return errors.New("unknown storage type")
	}

	if err != nil {
		return errors.Wrap(err, "could not connect to datastore")
	}

	if isFirstRun {
		if err = s.SeedUsers(ctx, datastoreRepository, authentication); err != nil {
			return errors.Wrap(err, "could not seed users")
		}
	}

	auditLog := service.CreateService(datastoreRepository, authentication, s.config.ServiceConfig)

	s.datastore = datastoreRepository
	s.authentication = authentication
	s.service = auditLog
	s.restHandler = rest.CreateHandler(auditLog)
	return nil
}

// SeedUsers initialize users database (it's just a quick and dirty shortcut)
func (s *server) SeedUsers(ctx context.Context, datastoreRepository service.DatastoreRepository, authentication auth.Auth) error {
	users := []struct {
		username       string
		password       string
		hasReadAccess  bool
		hasWriteAccess bool
	}{
		// full access user
		{
			username:       "admin",
			password:       "admin",
			hasReadAccess:  true,
			hasWriteAccess: true,
		},
		// read-only user
		{
			username:       "reader",
			password:       "reader",
			hasReadAccess:  true,
			hasWriteAccess: false,
		},
		// write-only user
		{
			username:       "writer",
			password:       "writer",
			hasReadAccess:  false,
			hasWriteAccess: true,
		},
	}

	for _, user := range users {
		hashedPassword, salt, err := authentication.EncryptPassword(user.password)
		if err != nil {
			return err
		}
		err = datastoreRepository.AddUser(ctx, user.username, hashedPassword, salt, user.hasReadAccess, user.hasWriteAccess)
		if err != nil {
			return err
		}
	}

	return nil
}

// Start starts the application in blocking mode
func (s *server) Start(ctx context.Context, addr string) {
	// Create Router for HTTP Server
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.SetupRouter(ctx, s.restHandler),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	// Start background workers
	go s.StartBackgroundWorkers(ctx)

	// Start REST Server in Blocking mode
	log.Printf("[OK] Starting HTTP Server on %s\n", addr)
	err := s.httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalln(err.Error())
	}

	// Code Reach Here after HTTP Server Shutdown!
	log.Println("[OK] HTTP REST Server is shutting down!")
}

// StartBackgroundWorkers starts background worker of services
func (s *server) StartBackgroundWorkers(ctx context.Context) {
	go s.service.StartBackgroundWorkers(ctx)
	<-ctx.Done()
}

// GracefulShutdown listen over the quitSignal to graceful shutdown the app
func (s *server) GracefulShutdown(quitSignal <-chan os.Signal, done chan<- bool) {
	// Wait for OS signals
	<-quitSignal

	// Create a 5s timeout context or waiting for app to shut down after 5 seconds
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()

	s.httpServer.SetKeepAlivesEnabled(false)
	if err := s.httpServer.Shutdown(ctxTimeout); err != nil {
		log.Println(err.Error())
	}
	log.Println("HTTP REST Server graceful shutdown completed")

	if s.datastore != nil {
		err := s.datastore.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}

	close(done)
}
