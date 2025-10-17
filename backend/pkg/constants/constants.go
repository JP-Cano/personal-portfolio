package constants

import "errors"

var (
	ErrInvalidInput       = errors.New("invalid input data")
	ErrExperienceNotFound = errors.New("experience not found")
	ErrProjectNotFound    = errors.New("project not found")
)
