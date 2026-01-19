package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mahimtalukder/go-student-api/internal/config"
	"github.com/mahimtalukder/go-student-api/internal/http/handlers/student"
)

func main() {
	// Load configuration (env/file) and fail fast if required values are missing.
	cfg := config.MustLoad()

	// Set up the HTTP router and register endpoints.
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New())

	// Create the HTTP server with the configured address and router.
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}
	slog.Info("Server is listening on " + cfg.Address)

	// Prepare for graceful shutdown by listening for OS interrupt/terminate signals.
	// When a signal arrives, we stop accepting new requests and let in-flight ones finish.
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine so main can wait for a shutdown signal.
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Stoping the server")
		}
	}()

	// Block until we receive a shutdown signal (graceful shutdown trigger).
	<-done

	// Allow a short window for in-flight requests to finish before forcing exit.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	slog.Info("Shutting down server")
	// Graceful shutdown: stop new connections and wait until ctx deadline.
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("Server shut down")
}
