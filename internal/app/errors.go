package app

import (
	"errors"
	"fmt"
)

var ErrProfileNotFound = errors.New("profile not found")

type ProfileNotFoundError struct {
	Name string
}

func (e *ProfileNotFoundError) Error() string {
	return fmt.Sprintf("%v: %s", ErrProfileNotFound, e.Name)
}

func (e *ProfileNotFoundError) Unwrap() error {
	return ErrProfileNotFound
}
