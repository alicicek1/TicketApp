package util

import "net/http"

type Error struct {
	ApplicationName string `json:"applicationName"`
	Operation       string `json:"operation"`
	Description     string `json:"description"`
	StatusCode      int    `json:"statusCode"`
	ErrorCode       int    `json:"errorCode "`
}

func NewError(applicationName, operation, description string, statusCode, errorCode int) *Error {
	return &Error{
		ApplicationName: applicationName,
		Operation:       operation,
		Description:     description,
		StatusCode:      statusCode,
		ErrorCode:       errorCode,
	}
}

func (e *Error) ModifyDescription(desc string) *Error {
	e.Description = desc
	return e
}

func (e *Error) ModifyErrorCode(code int) *Error {
	e.StatusCode = code
	return e
}

var (
	UnKnownError = NewError("-", "-", "An unknown error occurred.", 1, -1)
	NotFound     = NewError("-", "GET", "Not found.", http.StatusNotFound, -1)
)
