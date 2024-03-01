package model

import "errors"

var (
	ErrNoFile                 = errors.New("no file")
	ErrToManyFiles            = errors.New("too many files")
	ErrFileExtensionViolation = errors.New("file extension violation")
	ErrInvalidFileContent     = errors.New("invalid file content")
)
