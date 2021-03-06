package exception

import "fmt"

const (
	InvalidParameters int = 400

	Conflict int = 409

	NotFound int = 404

	ProcessmentError int = 500

	InvalidCredentials int = 401

	Unknown int = 500
)

// Error is a struct for error
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Err     error  `json:"-"`
}

// New creates a new error
func New(code int, message string, err error) error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("code %d , message %s, %q", e.Code, e.Message, e.Err)
}
