package weather

import "errors"

var (
	// ErrInvalidZipcode is returned when the CEP is not exactly 8 digits.
	ErrInvalidZipcode = errors.New("invalid zipcode")
	// ErrZipcodeNotFound is returned when the CEP is well formed but unknown.
	ErrZipcodeNotFound = errors.New("can not find zipcode")
	// ErrUpstream is returned when a location or weather provider fails.
	ErrUpstream = errors.New("upstream request failed")
)
