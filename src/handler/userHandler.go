package handler

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
	"ticketApp/src/config"
	"ticketApp/src/service"
	"ticketApp/src/type/entity"
	"ticketApp/src/type/util"
)

type UserHandler struct {
	userService service.UserService
	cfg         *config.AppConfig
}

func NewUserHandler(userService service.UserService, cfg *config.AppConfig) UserHandler {
	return UserHandler{userService: userService, cfg: cfg}
}

func (h *UserHandler) UserGetById(ctx echo.Context) error {
	id := ctx.Param("id")

	if id == "" {
		return ctx.JSON(http.StatusBadRequest, util.PathVariableNotFound.ModifyApplicationName("user handler").ModifyErrorCode(4008))
	}

	if util.IsValidUUID(id) {
		return ctx.JSON(http.StatusBadRequest, util.PathVariableIsNotValid.ModifyApplicationName("user handler").ModifyErrorCode(4009))
	}

	user, errSrv := h.userService.UserServiceGetById(id)
	if errSrv != nil || user == nil {
		return ctx.JSON(http.StatusNotFound, util.NotFound.ModifyApplicationName("user handler").ModifyErrorCode(4010))
	}

	return ctx.JSON(http.StatusOK, user)
}
func (h *UserHandler) UserUpsert(ctx echo.Context) error {
	userPostRequestModel := entity.UserPostRequestModel{}

	err := json.NewDecoder(ctx.Request().Body).Decode(&userPostRequestModel)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, util.InvalidBody.ModifyApplicationName("user handler").ModifyErrorCode(4012))
	}
	user := entity.User{Username: userPostRequestModel.Username,
		Password: userPostRequestModel.Password,
		Email:    userPostRequestModel.Email,
		Type:     userPostRequestModel.Type,
		Age:      userPostRequestModel.Age,
	}

	id := ctx.QueryParam("id")
	if id != "" {
		if util.IsValidUUID(id) {
			return ctx.JSON(http.StatusBadRequest, util.NewError("user handler", "POST", "Provided identifier is not valid format.", http.StatusBadRequest, 4011))
		}
		user.Id = id
	}

	res, errSrv := h.userService.UserServiceInsert(user)
	if errSrv != nil {
		return ctx.JSON(errSrv.StatusCode, errSrv)
	}
	return ctx.JSON(http.StatusOK, res)
}
func (h *UserHandler) UserDeleteById(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, util.PathVariableNotFound.ModifyApplicationName("user handler").ModifyErrorCode(4013))
	}

	if util.IsValidUUID(id) {
		return ctx.JSON(http.StatusBadRequest, util.PathVariableIsNotValid.ModifyApplicationName("user handler").ModifyErrorCode(4014))
	}

	res, errSrv := h.userService.UserServiceDeleteById(id)
	if errSrv != nil {
		return ctx.JSON(errSrv.ErrorCode, util.NewError(errSrv.ApplicationName, errSrv.Operation, errSrv.Description, errSrv.ErrorCode, errSrv.StatusCode))
	}

	return ctx.JSON(http.StatusOK, res)
}
func (h *UserHandler) UserGetAll(ctx echo.Context) error {
	filter := util.Filter{}
	page, pageSize := util.ValidatePaginationFilters(ctx.QueryParam("page"), ctx.QueryParam("pageSize"), h.cfg.MaxPageLimit)
	filter.Page = page
	filter.PageSize = pageSize

	sortingField, sortingDirection := util.ValidateSortingFilters(entity.Category{}, ctx.QueryParam("sort"), ctx.QueryParam("sDirection"))
	filter.SortingField = sortingField
	filter.SortingDirection = sortingDirection

	filters := map[string]interface{}{}
	if username := ctx.QueryParam("username"); username != "" && len(username) < 30 {
		filters["username"] = bson.M{"$regex": primitive.Regex{
			Pattern: username,
			Options: "i",
		}}
	}

	if mingAgeStr := ctx.QueryParam("minAge"); mingAgeStr != "" {
		if minAge, err := strconv.Atoi(mingAgeStr); err == nil {
			filters["age"] = bson.M{"$gte": minAge}
		}
	}

	if maxAgeStr := ctx.QueryParam("maxAge"); maxAgeStr != "" {
		if maxAge, err := strconv.Atoi(maxAgeStr); err == nil {
			minFilter, exist := filters["age"]
			if exist {
				delete(filters, "age")
				filters["$and"] = bson.A{
					bson.M{"age": minFilter},
					bson.M{"age": bson.M{"$lte": maxAge}},
				}
			} else {
				filters["age"] = bson.M{"$lte": maxAge}
			}
		}
	}

	filter.Filters = filters

	res, err := h.userService.UserServiceGetAll(filter)
	if err != nil {
		return ctx.JSON(err.StatusCode, err)
	}

	if res.Users == nil {
		return ctx.JSON(http.StatusNotFound, util.NotFound.ModifyApplicationName("user handler").ModifyErrorCode(5001))
	}
	ctx.Response().Header().Add("x-total-count", strconv.FormatInt(res.RowCount, 10))
	return ctx.JSON(http.StatusOK, res)
}
