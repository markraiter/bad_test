package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/markraiter/bad_test/docs"

	"github.com/go-playground/validator"
	"github.com/markraiter/bad_test/internal/app/api"
	"github.com/markraiter/bad_test/internal/app/api/handler"
	"github.com/markraiter/bad_test/internal/app/service"
	"github.com/markraiter/bad_test/internal/config"
)

const (
	timoutLimit = 5
)

// @title BAD test API
// @version	1.0
// @description	This is a test API for BAD test.
// @contact.name Mark Raiter
// @contact.email raitermark@proton.me
// @host localhost:5555
// @BasePath /api/v1
func main() {
	// Initialize config.
	cfg := config.MustLoad()

	// Initialize logger.
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// Initialize validator.
	// Register custom validators.
	validate := validator.New()

	log.Info("Starting application...")
	log.Info("port: " + cfg.Server.AppAddress)

	// Initialize service layer.
	service := service.New(log)

	// Initialize transport layer.
	handler := handler.New(service, log, validate, cfg)

	// Initialize server.
	server := api.New(cfg, handler)

	// Initialize stop channel for graceful shutdown.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Starting listening port concurrently.
	go func() {
		if err := server.HTTPServer.Listen(cfg.Server.AppAddress); err != nil {
			log.Error("HTTPServer.Listen", err)
		}
	}()

	<-stop

	// Gracefully stop.

	if err := server.HTTPServer.ShutdownWithTimeout(timoutLimit * time.Second); err != nil {
		log.Error("ShutdownWithTimeout", err)
	}

	if err := server.HTTPServer.Shutdown(); err != nil {
		log.Error("Shutdown", err)
	}

	log.Info("server stopped")
}
