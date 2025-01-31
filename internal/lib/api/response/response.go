package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOk,
	}
}

func OKWithData(data interface{}) Response {
	return Response{
		Status: StatusOk,
		Data:   data,
	}
}

func Error(msg string) Response {
	return Response{
		Status:  StatusError,
		Message: msg,
	}
}

func InternalError() Response {
	return Error("iternal error")
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field % s is not valid", err.Field()))
		}
	}

	return Response{
		Status:  StatusError,
		Message: strings.Join(errMsgs, ", "),
	}
}
