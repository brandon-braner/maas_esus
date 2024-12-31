package errors

import "fmt"

type CustomError struct {
	ErrorMessage string `json:"error_message"`
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("errorMessage: %s", e.ErrorMessage)
}
