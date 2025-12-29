package utils

import (
	"context"
	"errors"
	"time"
)

// Retry2 is a helper func to retry a function fc a specified number of times if it encounters an error.
func Retry2[T any](fc func() (T, error), duration time.Duration, times int) (T, error) {
	var result T
	var err error

	for range times {
		result, err = fc()
		if err == nil {
			return result, err
		}
		time.Sleep(duration)
	}
	return result, err
}

// Retry2WithContext is a helper func to retry a function fc a specified number of times if it encounters an error.
func Retry2WithContext[T any](ctx context.Context, fc func() (T, error), duration time.Duration, times int) (T, error) {
	var result T
	var err error

	for range times {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
			result, err = fc()
			if err == nil {
				return result, err
			}
			time.Sleep(duration)
		}
	}
	return result, err
}

// Retry2OnError is a helper func to retry a function fc a specified number of times if it encounters an error.
// If the error is not the same as errRetry, it will return the error.
func Retry2OnError[T any](ctx context.Context, fc func() (T, error), errRetry error, duration time.Duration, times int) (T, error) {
	var result T
	var err error

	for range times {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
			result, err = fc()
			if err == nil {
				return result, err
			}
			if !errors.Is(err, errRetry) {
				return result, err
			}
			time.Sleep(duration)
		}
	}
	return result, err
}

// Retry is a helper func to retry a function fc a specified number of times if it encounters an error.
func Retry(fc func() error, duration time.Duration, times int) error {
	var err error
	for range times {
		err = fc()
		if err == nil {
			return err
		}
		time.Sleep(duration)
	}
	return err
}

// RetryWithContext is a helper func to retry a function fc a specified number of times if it encounters an error.
func RetryWithContext(ctx context.Context, fc func() error, duration time.Duration, times int) error {
	var err error
	for range times {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err = fc()
			if err == nil {
				return err
			}
			time.Sleep(duration)
		}
	}
	return err
}

func RetryOnError(ctx context.Context, fc func() error, errRetry error, duration time.Duration, times int) error {
	var err error
	for range times {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err = fc()
			if err == nil {
				return err
			}
			if !errors.Is(err, errRetry) {
				return err
			}

			time.Sleep(duration)
		}
	}
	return err
}
