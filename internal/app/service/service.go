package service

import "log/slog"

type Services struct {
	TaskService
}

// New returns new instance of the Service.
func New(log *slog.Logger) *Services {
	return &Services{
		TaskService{log: log},
	}
}
