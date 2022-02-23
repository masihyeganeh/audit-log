package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/masihyeganeh/audit-log/internal/service"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	GitCommit = "Insider"
	BuildTime string
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func main() {
	// Setting up the main context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Print Start Ascii Art
	printAsciiArt()

	// Configuration
	// TODO: Move to a config file
	cfg := Config{
		ServerAddress:               ":80",
		DataStorage:                 "clickhouse_with_map", // clickhouse_with_map | clickhouse_with_nested
		DataStorageConnectionString: "tcp://datastore:9000",
		ServiceConfig: service.Config{
			ChannelSize:             10_000,
			MaxWorkers:              10,
			EventsBufferSize:        10,
			MaxEventsBufferDuration: 1 * time.Second,
		},
		JwtSecret: "jwt_secret_key_here",
	}

	// Create New Server
	server := NewServer(&cfg)

	// Initialize the Server Dependencies
	err := server.Initialize(ctx)
	if err != nil {
		log.Fatalf("failed to initialize server : %v", err)
	}

	done := make(chan bool, 1)
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, syscall.SIGINT, syscall.SIGTERM)

	// Graceful shutdown goroutine
	go server.GracefulShutdown(quitSignal, done)

	// Start server in blocking mode
	server.Start(ctx, cfg.ServerAddress)

	// Wait for graceful shutdown signal
	<-done

	// Kill other background jobs
	cancel()
	log.Println("Waiting for background jobs to finish their works...")

	// Wait for all other background jobs to finish their works
	server.Wait()

	log.Println("App Shutdown successfully")
}

func printAsciiArt() {
	if len(BuildTime) == 0 {
		BuildTime = time.Now().Format("2006-01-02-15:04:05")
	}

	asciiArt := `
╔════════════════════════════════════════════════════════════════════════════════════╗
║         █████  ██    ██ ██████  ██ ████████       ██       ██████   ██████         ║
║        ██   ██ ██    ██ ██   ██ ██    ██          ██      ██    ██ ██              ║
║        ███████ ██    ██ ██   ██ ██    ██    █████ ██      ██    ██ ██   ███        ║
║        ██   ██ ██    ██ ██   ██ ██    ██          ██      ██    ██ ██    ██        ║
║        ██   ██  ██████  ██████  ██    ██          ███████  ██████   ██████         ║
║                                                                                    ║
║             Version: VERSION           Build Time: TIME OF COMPILATION             ║
╚════════════════════════════════════════════════════════════════════════════════════╝
`
	asciiArt = strings.ReplaceAll(asciiArt, "VERSION", GitCommit[0:7])
	asciiArt = strings.ReplaceAll(asciiArt, "TIME OF COMPILATION", BuildTime)
	fmt.Println(asciiArt)
}
