package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

// ValidationError constructs a Response with a Status of "Error" and an Error containing
// a comma-separated list of human-readable error messages for the given validator.ValidationErrors.
func ValidationError(errors validator.ValidationErrors) Response {
	var errorMessages []string

	for _, err := range errors {
		field := err.Field()
		tag := err.ActualTag()

		switch tag {
		case "required":
			errorMessages = append(errorMessages, fmt.Sprintf("field %s is a required field", field))
		case "url":
			errorMessages = append(errorMessages, fmt.Sprintf("field %s is not a valid URL", field))
		default:
			errorMessages = append(errorMessages, fmt.Sprintf("field %s is not valid", field))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errorMessages, ", "),
	}
}
