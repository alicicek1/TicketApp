package handler

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"ticketApp/src/config"
	"ticketApp/src/service"
	"ticketApp/src/type/entity"
	"ticketApp/src/type/util"
)

type CategoryHandler struct {
	categoryService service.CategoryService
	cfg             *config.AppConfig
}

func NewCategoryHandler(categoryService service.CategoryService, cfg *config.AppConfig) CategoryHandler {
	return CategoryHandler{
		categoryService: categoryService,
		cfg:             cfg,
	}
}

func (h *CategoryHandler) CategoryGetById(ctx echo.Context) error {
	id := ctx.Param("id")

	if id == "" {
		return ctx.JSON(http.StatusBadRequest, util.PathVariableNotFound.ModifyApplicationName("category handler").ModifyErrorCode(4016))
	}

	if util.IsValidUUID(id) {
		return ctx.JSON(http.StatusBadRequest, util.PathVariableIsNotValid.ModifyApplicationName("category handler").ModifyErrorCode(4017))
	}

	category, errSrv := h.categoryService.CategoryServiceGetById(id)
	if errSrv != nil || category == nil {
		return ctx.JSON(http.StatusNotFound, util.NotFound.ModifyApplicationName("category handler").ModifyErrorCode(4018))
	}

	return ctx.JSON(http.StatusOK, category)
}
func (h *CategoryHandler) CategoryInsert(ctx echo.Context) error {
	categoryPostRequestModel := entity.CategoryPostRequestModel{}

	err := json.NewDecoder(ctx.Request().Body).Decode(&categoryPostRequestModel)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, util.InvalidBody.ModifyApplicationName("category handler").ModifyErrorCode(4022))
	}
	category := entity.Category{
		Name: categoryPostRequestModel.Name,
	}

	res, errSrv := h.categoryService.CategoryServiceInsert(category)
	if errSrv != nil {
		return ctx.JSON(errSrv.ErrorCode, util.NewError(errSrv.ApplicationName, errSrv.Operation, errSrv.Description, errSrv.ErrorCode, errSrv.StatusCode))
	}
	return ctx.JSON(http.StatusOK, res)
}
func (h *CategoryHandler) CategoryDeleteById(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, util.PathVariableNotFound.ModifyApplicationName("user handler").ModifyErrorCode(4013))
	}

	if util.IsValidUUID(id) {
		return ctx.JSON(http.StatusBadRequest, util.PathVariableIsNotValid.ModifyApplicationName("user handler").ModifyErrorCode(4014))
	}

	res, errSrv := h.categoryService.CategoryServiceDeleteById(id)
	if errSrv != nil {
		return ctx.JSON(errSrv.ErrorCode, util.NewError(errSrv.ApplicationName, errSrv.Operation, errSrv.Description, errSrv.ErrorCode, errSrv.StatusCode))
	}

	return ctx.JSON(http.StatusOK, res)
}
func (h *CategoryHandler) CategoryGetAll(ctx echo.Context) error {
	filter := util.Filter{}
	page, pageSize := util.ValidatePaginationFilters(ctx.QueryParam("page"), ctx.QueryParam("pageSize"), h.cfg.MaxPageLimit)
	filter.Page = page
	filter.PageSize = pageSize

	sortingField, sortingDirection := util.ValidateSortingFilters(entity.Category{}, ctx.QueryParam("sort"), ctx.QueryParam("sDirection"))
	filter.SortingField = sortingField
	filter.SortingDirection = sortingDirection

	nameValue := ctx.QueryParam("name")
	if nameValue != "" {
		filter.Filters["name"] = util.CreateEqualFilter(nameValue, "name")
	}

	categories, errSrv := h.categoryService.CategoryServiceGetAll(filter)
	if errSrv != nil {
		return ctx.JSON(errSrv.StatusCode, errSrv)
	}

	if categories == nil {
		return ctx.JSON(http.StatusNotFound, util.NotFound.ModifyApplicationName("category handler").ModifyErrorCode(5000))
	}
	return ctx.JSON(http.StatusOK, categories)
}
