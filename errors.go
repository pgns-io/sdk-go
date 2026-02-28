package sdk

import (
	"errors"
	"fmt"
)

// PigeonsError is returned when the pgns API responds with a non-2xx status.
type PigeonsError struct {
	Message    string
	StatusCode int
}

func (e *PigeonsError) Error() string {
	return fmt.Sprintf("pgns: %s (status %d)", e.Message, e.StatusCode)
}

// IsNotFound reports whether err is a [PigeonsError] with status 404.
func IsNotFound(err error) bool {
	var pe *PigeonsError
	return errors.As(err, &pe) && pe.StatusCode == 404
}

// IsUnauthorized reports whether err is a [PigeonsError] with status 401.
func IsUnauthorized(err error) bool {
	var pe *PigeonsError
	return errors.As(err, &pe) && pe.StatusCode == 401
}
