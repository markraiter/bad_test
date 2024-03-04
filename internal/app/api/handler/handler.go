package handler

import (
	"log/slog"
)

type ServiceInterface interface {
	Processor
}

type Handler struct {
	TaskHandler
}

// New returns new instance of the Handler.
func New(services ServiceInterface, log *slog.Logger) *Handler {
	return &Handler{
		TaskHandler{log: log, processor: services},
	}
}
