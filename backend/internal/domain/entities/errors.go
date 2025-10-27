package entities

import "errors"

var (
	ErrGitHubTokenNotFound = errors.New("GitHub token not found, please configure your GitHub token in settings")
	ErrUnauthorized        = errors.New("unauthorized access")
	ErrNotFound            = errors.New("resource not found")
	ErrInvalidInput        = errors.New("invalid input")
	ErrInternalServer      = errors.New("internal server error")
)

