package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"ticketApp/src/service"
	"ticketApp/src/type/entity"
	"ticketApp/src/type/util"
	"ticketApp/src/type/util/request"
	"ticketApp/src/type/util/response"
	"time"
)

type handler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) handler {
	return handler{userService}
}

func (h *handler) UserGetById(ctx echo.Context) error {
	customEchoCtx := response.NewCustomEchoContext(&ctx)
	id := ctx.Param("id")
	if id == "" {
		return customEchoCtx.ReturnBadRequestResponse(*util.NewError("", "", "Path variable not found.", http.StatusBadRequest, util.PATH_VARIABLE_NOT_FOUND_ERROR_CODE))
	}

	user, err := h.userService.UserServiceGetById(id)
	if err != nil {
		return customEchoCtx.ReturnBadRequestResponse(*util.NewError("", "", err.Error(), http.StatusNotFound, util.USER_GET_BY_ID_ERROR_CODE))
	}

	if user == nil {
		return customEchoCtx.ReturnBadRequestResponse(*util.NewError("", "", "User not found.", http.StatusNotFound, util.USER_NOT_FOUND_ERROR_CODE))
	}

	return customEchoCtx.ReturnOkResponseWithBody(user)
}

func (h *handler) UserUpsert(ctx echo.Context) error {
	customEchoCtx := response.NewCustomEchoContext(&ctx)
	userPostRequestModel := request.UserPostRequestModel{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&userPostRequestModel)
	if err != nil {
		return customEchoCtx.ReturnBadRequestResponse(*util.NewError("", "", err.Error(), http.StatusBadRequest, util.USER_INVALID_POST_BODY_ERROR_CODE))
	}

	empId := uuid.New().String()
	user := entity.User{
		BaseEntity: entity.BaseEntity{
			Id:        empId,
			CreatedAt: time.Now(),
		},
		Username: userPostRequestModel.Username,
		Password: userPostRequestModel.Password,
		Email:    userPostRequestModel.Email,
		Type:     userPostRequestModel.Type,
	}
	res, err := h.userService.UserServiceInsert(user)
	if err != nil {
		return customEchoCtx.ReturnBadRequestResponse(*util.NewError("", "", err.Error(), http.StatusBadRequest, util.USER_VALIDATION_ERROR_CODE))
	}

	return customEchoCtx.ReturnCreatedResponseWithBody(res)
}

func (h *handler) UserDeleteById(ctx echo.Context) error {
	customEchoCtx := response.NewCustomEchoContext(&ctx)
	id := ctx.Param("id")
	if id == "" {
		return customEchoCtx.ReturnBadRequestResponse(*util.NewError("", "", "Path variable not found.", http.StatusBadRequest, util.PATH_VARIABLE_NOT_FOUND_ERROR_CODE))
	}

	res, err := h.userService.UserServiceDeleteById(id)
	if err != nil {
		return customEchoCtx.ReturnBadRequestResponse(*util.NewError("", "", err.Error(), http.StatusBadRequest, util.USER_DELETE_BY_ID_ERROR_CODE))
	}

	return customEchoCtx.ReturnOkResponseWithBody(res)
}
