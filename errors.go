package ccatapi

import "fmt"

var (
	ErrUploadMissingFile = fmt.Errorf("missing file, cannot upload")
)

// errUnknownError creates an unknown error with the given status code and response text.
func errUnknownError(statusCode int, respText string) error {
	return fmt.Errorf("unknown error: %d - %s", statusCode, respText)
}
