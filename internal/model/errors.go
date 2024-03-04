package model

import (
	"errors"
	"log/slog"
)

var (
	ErrNoFile                 = errors.New("no file")
	ErrToManyFiles            = errors.New("too many files")
	ErrFileExtensionViolation = errors.New("file extension violation")
	ErrInvalidFileContent     = errors.New("invalid file content")
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
