package service

import "log/slog"

type Services struct {
}

// New returns new instance of the Service.
func New(log *slog.Logger) *Services {
	return &Services{}
}
