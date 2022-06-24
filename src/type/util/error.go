package util

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

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
	e.ErrorCode = code
	return e
}

func (e *Error) ModifyApplicationName(application string) *Error {
	e.ApplicationName = application
	return e
}

func (e *Error) ModifyOperation(operation string) *Error {
	e.Operation = operation
	return e
}

var (
	UnKnownError           = NewError("-", "-", "An unknown error occurred.", 1, -1)
	PathVariableNotFound   = NewError("", "GET", "Path variable not found.", http.StatusBadRequest, -1)
	PathVariableIsNotValid = NewError("", "GET", "Path variable is not valid format.", http.StatusBadRequest, -1)
	NotFound               = NewError("", "GET", "Not found.", http.StatusNotFound, -1)
	InvalidBody            = NewError("", "UPSERT", "Request body is not valid.", http.StatusBadRequest, -1)
	UpsertFailed           = NewError("", "UPSERT", "Upsert failed.", http.StatusBadRequest, -1)
	PostValidation         = NewError("", "POST", "", http.StatusBadRequest, -1)
	DeleteFailed           = NewError("", "DELETE", "failed delete", http.StatusBadRequest, -1)
	CountGet               = NewError("", "GET", "count get failed", http.StatusBadRequest, -1)
)

type (
	httpErrorHandler struct {
		statusCodes map[error]int
	}
)

func NewHttpErrorHandler(errorStatusCodeMaps map[error]int) *httpErrorHandler {
	return &httpErrorHandler{
		statusCodes: errorStatusCodeMaps,
	}
}

func (self *httpErrorHandler) getStatusCode(err error) int {
	for key, value := range self.statusCodes {
		if errors.Is(err, key) {
			return value
		}
	}

	return http.StatusInternalServerError
}

func unwrapRecursive(err error) error {
	var originalErr = err

	for originalErr != nil {
		var internalErr = errors.Unwrap(originalErr)

		if internalErr == nil {
			break
		}

		originalErr = internalErr
	}

	return originalErr
}

func (self *httpErrorHandler) Handler(err error, c echo.Context) {
	he, ok := err.(*echo.HTTPError)
	if ok {
		if he.Internal != nil {
			if herr, ok := he.Internal.(*echo.HTTPError); ok {
				he = herr
			}
		}
	} else {
		he = &echo.HTTPError{
			Code:    self.getStatusCode(err),
			Message: unwrapRecursive(err).Error(),
		}
	}

	code := he.Code
	message := he.Message
	if _, ok := he.Message.(string); ok {
		message = map[string]interface{}{"message": err.Error()}
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(he.Code)
		} else {
			err = c.JSON(code, message)
		}
		if err != nil {
			c.Echo().Logger.Error(err)
		}
	}
}

var ErrDocumentNotFound = errors.New("DocumentNotFound")

func NewErrorStatusCodeMaps() map[error]int {

	var errorStatusCodeMaps = make(map[error]int)
	errorStatusCodeMaps[ErrDocumentNotFound] = http.StatusNotFound
	return errorStatusCodeMaps
}
