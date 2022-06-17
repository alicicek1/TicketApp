package response

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ticketApp/src/type/util"
)

type CustomEchoContext struct {
	context echo.Context
}

func NewCustomEchoContext(e *echo.Context) CustomEchoContext {
	return CustomEchoContext{context: *e}
}

func (c CustomEchoContext) ReturnBadRequestResponse(err util.Error) error {
	return c.context.JSON(http.StatusBadRequest, err)
}

func (c CustomEchoContext) ReturnOkResponseWithBody(i interface{}) error {
	return c.context.JSON(http.StatusOK, i)
}

func (c CustomEchoContext) ReturnCreatedResponseWithBody(i interface{}) error {
	return c.context.JSON(http.StatusCreated, i)
}

func (c CustomEchoContext) ReturnNotFoundResponse() error {
	return c.context.JSON(http.StatusNotFound, "Not found.")
}

func (c CustomEchoContext) ReturnNotFoundUpdateResponse() error {
	return c.context.JSON(http.StatusNotFound, "There is no model by provided identifier.")
}
