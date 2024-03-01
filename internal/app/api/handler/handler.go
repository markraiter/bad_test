package handler

import (
	"log/slog"

	"github.com/go-playground/validator"
	"github.com/markraiter/bad_test/internal/config"
)

type ServiceInterface interface {
}

type Handler struct {
	HealthCheck
}

// New returns new instance of the Handler.
func New(services ServiceInterface, log *slog.Logger, validate *validator.Validate, cfg *config.Config) *Handler {
	return &Handler{
		HealthCheck{log: log},
	}
}
