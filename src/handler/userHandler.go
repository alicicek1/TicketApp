package handler

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"reflect"
	"ticketApp/src/service"
	"ticketApp/src/type/entity"
	"ticketApp/src/type/util"
	"time"
)

type handler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) handler {
	return handler{userService}
}

func (h *handler) UserGetById(ctx echo.Context) error {
	id := ctx.QueryParam("id")
	if id == "" {
		id = ctx.Param("id")
	}
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, *util.NewError("user handler", "GET", "Path variable not found.", http.StatusBadRequest, 4008))
	}

	if reflect.TypeOf(id) != reflect.TypeOf("string") {
		return ctx.JSON(http.StatusBadRequest, *util.NewError("user handler", "GET", "Path variable is not valid format.", http.StatusBadRequest, 4009))
	}

	user, errSrv := h.userService.UserServiceGetById(id)
	if errSrv != nil || user == nil {
		return ctx.JSON(http.StatusNotFound, *util.NewError("user handler", "GET", "User not found.", http.StatusNotFound, 4010))
	}

	return ctx.JSON(http.StatusOK, user)
}

func (h *handler) UserUpsert(ctx echo.Context) error {
	userPostRequestModel := util.UserPostRequestModel{}

	err := json.NewDecoder(ctx.Request().Body).Decode(&userPostRequestModel)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, *util.NewError("user handler", "UPSERT", err.Error(), http.StatusBadRequest, 4012))
	}
	user := entity.User{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Username:  userPostRequestModel.Username,
		Password:  userPostRequestModel.Password,
		Email:     userPostRequestModel.Email,
		Type:      userPostRequestModel.Type,
	}

	var empId string
	id := ctx.QueryParam("id")
	if id != "" {
		if reflect.TypeOf(id) != reflect.TypeOf("string") {
			return ctx.JSON(http.StatusBadRequest, *util.NewError("user handler", "POST", "Provided identifier is not valid format.", http.StatusBadRequest, 4011))
		}
		empId = id

		empIdStr, err := primitive.ObjectIDFromHex(empId)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, *util.NewError("user handler", "POST", "Provided identifier is not valid format.", http.StatusBadRequest, 4012))
		}
		user.Id = &empIdStr
	}

	res, errSrv := h.userService.UserServiceInsert(user)
	if errSrv != nil {
		return ctx.JSON(errSrv.ErrorCode, *util.NewError(errSrv.ApplicationName, errSrv.Operation, errSrv.Description, errSrv.ErrorCode, errSrv.StatusCode))
	}
	return ctx.JSON(http.StatusOK, res)
}

func (h *handler) UserDeleteById(ctx echo.Context) error {
	id := ctx.QueryParam("id")
	if id == "" {
		id = ctx.Param("id")
	}
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, *util.NewError("user handler", "DELETE", "Path variable not found.", http.StatusBadRequest, 4013))
	}

	if reflect.TypeOf(id) != reflect.TypeOf("string") {
		return ctx.JSON(http.StatusBadRequest, *util.NewError("user handler", "GET", "Path variable is not valid format.", http.StatusBadRequest, 4014))
	}

	res, errSrv := h.userService.UserServiceDeleteById(id)
	if errSrv != nil {
		return ctx.JSON(errSrv.ErrorCode, *util.NewError(errSrv.ApplicationName, errSrv.Operation, errSrv.Description, errSrv.ErrorCode, errSrv.StatusCode))
	}

	return ctx.JSON(http.StatusOK, util.UserDeleteResponseType{
		IsSuccess: res,
	})
}
