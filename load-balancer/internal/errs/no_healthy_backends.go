package errs

import "errors"

var (
	ErrNoHealthyBackends = errors.New("no healthy backends available")
)
