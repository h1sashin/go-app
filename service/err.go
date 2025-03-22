package service

import (
	"errors"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrBanned       = errors.New("banned")
	ErrConflict     = errors.New("conflict")
	ErrInternal     = errors.New("internal error")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrExpired      = errors.New("expired")
	ErrInvalidToken = errors.New("invalid token")
	ErrInvalidRole  = errors.New("invalid role")
)
