package domain

import (
	"errors"
	"fmt"
)

var (
	// ErrInternal is an error for when an internal service fails
	ErrInternal = errors.New("internal error")
	// ErrDataNotFound is an error for when requested data is not found
	ErrDataNotFound = errors.New("data not found")
	// ErrNoUpdatedData is an error for when no data is provided to update
	ErrNoUpdatedData = errors.New("no data to update")
	// ErrDataConflict is an error for data conflict
	ErrDataConflict = errors.New("data conflict error")
	// ErrConflictingData is an error for when data conflicts with existing data
	ErrConflictingData = errors.New("data conflicts with existing data in unique column")
	// ErrForbidden is an error for when the user is forbidden to access the resource
	ErrForbidden = errors.New("user is forbidden to access the resource")

	// ErrTokenDuration is an error for when the token duration format is invalid
	ErrTokenDuration = errors.New("invalid token duration format")
	// ErrTokenCreation is an error for when the token creation fails
	ErrTokenCreation = errors.New("error creating token")
	// ErrExpiredToken is an error for when the access token is expired
	ErrExpiredToken = errors.New("access token has expired")
	// ErrInvalidToken is an error for when the access token is invalid
	ErrInvalidToken = errors.New("access token is invalid")

	// ErrInvalidCredentials is an error for when the credentials are invalid
	ErrInvalidCredentials = errors.New("invalid email or password")
	// ErrEmptyAuthorizationHeader is an error for when the authorization header is empty
	ErrEmptyAuthorizationHeader = errors.New("authorization header is not provided")
	// ErrInvalidAuthorizationHeader is an error for when the authorization header is invalid
	ErrInvalidAuthorizationHeader = errors.New("authorization header format is invalid")
	// ErrInvalidAuthorizationType is an error for when the authorization type is invalid
	ErrInvalidAuthorizationType = errors.New("authorization type is not supported")
	// ErrUnauthorized is an error for when the user is unauthorized
	ErrUnauthorized = errors.New("user is unauthorized to access the resource")

	// ErrInvalidFileExt is an error for when file extension is invalid
	ErrInvalidFileExt = errors.New("file extension is invalid")
)

type AppError struct {
	// DomainError
	DomainError error
	// WrappedError
	WrappedError error `json:"-"`
}

func (e *AppError) Error() string {
	msg := e.DomainError.Error()
	if e.WrappedError != nil {
		return fmt.Sprintf("%s: %v", msg, e.WrappedError)
	}
	return msg
}

func (e *AppError) Is(target error) bool {
	return errors.Is(e.DomainError, target)
}

func (e *AppError) Unwrap() error {
	return e.WrappedError
}

func ErrW(domainErr error, actualErr error) error {
	return &AppError{
		DomainError:  domainErr,
		WrappedError: actualErr,
	}
}

func NewErrInternal(actualErr error) error {
	return &AppError{
		DomainError:  ErrInternal,
		WrappedError: actualErr,
	}
}
