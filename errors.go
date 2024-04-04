package ccatapi

import "fmt"

// ErrUnknownError creates an unknown error with the given status code and response text.
func ErrUnknownError(statusCode int, respText string) error {
	return fmt.Errorf("unknown error: %d - %s", statusCode, respText)
}
