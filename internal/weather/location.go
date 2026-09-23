package weather

import (
	"regexp"
	"strings"
)

// zipcodePattern matches a Brazilian CEP with exactly eight digits.
var zipcodePattern = regexp.MustCompile(`^[0-9]{8}$`)

// Location is the city resolved from a CEP.
type Location struct {
	City  string
	State string
}

// Query is the place string sent to the weather provider.
func (l Location) Query() string {
	city := strings.TrimSpace(l.City)
	state := strings.TrimSpace(l.State)
	if state == "" {
		return city + ", Brazil"
	}
	return city + ", " + state + ", Brazil"
}

// ValidateZipcode reports whether zipcode is an 8-digit CEP.
func ValidateZipcode(zipcode string) error {
	if !zipcodePattern.MatchString(zipcode) {
		return ErrInvalidZipcode
	}
	return nil
}
