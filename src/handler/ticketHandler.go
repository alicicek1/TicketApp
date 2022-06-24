package handler

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"ticketApp/src/config"
	"ticketApp/src/service"
	"ticketApp/src/type/entity"
	"ticketApp/src/type/util"
	"time"
)

type TicketHandler struct {
	ticketService service.TicketService
	cfg           *config.AppConfig
}

func NewTicketHandler(ticketService service.TicketService, cfg *config.AppConfig) TicketHandler {
	return TicketHandler{ticketService: ticketService, cfg: cfg}
}

func (h *TicketHandler) TicketGetById(ctx echo.Context) error {
	id := ctx.Param("id")

	if id == "" {
		return ctx.JSON(http.StatusBadRequest, util.PathVariableNotFound.ModifyApplicationName("ticket handler").ModifyErrorCode(4018))
	}

	if util.IsValidUUID(id) {
		return ctx.JSON(http.StatusBadRequest, util.PathVariableIsNotValid.ModifyApplicationName("ticket handler").ModifyErrorCode(4019))
	}

	ticket, errSrv := h.ticketService.TicketServiceGetById(id)
	if errSrv != nil || ticket == nil {
		return ctx.JSON(http.StatusNotFound, util.NotFound.ModifyApplicationName("category handler").ModifyErrorCode(4018))
	}

	return ctx.JSON(http.StatusOK, ticket)
}
func (h *TicketHandler) TicketInsert(ctx echo.Context) error {
	ticketPostRequestModel := entity.TicketPostRequestModel{}

	err := json.NewDecoder(ctx.Request().Body).Decode(&ticketPostRequestModel)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, util.InvalidBody.ModifyApplicationName("category handler").ModifyErrorCode(4022).ModifyOperation("POST"))
	}
	category := entity.Ticket{
		Category:       ticketPostRequestModel.Category,
		Attachments:    ticketPostRequestModel.Attachments,
		Answers:        ticketPostRequestModel.Answers,
		Subject:        ticketPostRequestModel.Subject,
		Body:           ticketPostRequestModel.Body,
		CreatedBy:      ticketPostRequestModel.CreatedBy,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		LastAnsweredAt: time.Now(),
		Status:         byte(entity.CREATED),
	}

	res, errSrv := h.ticketService.TicketServiceInsert(category)
	if errSrv != nil {
		return ctx.JSON(errSrv.ErrorCode, util.NewError(errSrv.ApplicationName, errSrv.Operation, errSrv.Description, errSrv.ErrorCode, errSrv.StatusCode))
	}
	return ctx.JSON(http.StatusOK, res)
}
func (h *TicketHandler) TicketDeleteById(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, util.PathVariableNotFound.ModifyApplicationName("user handler").ModifyErrorCode(4020))
	}

	if util.IsValidUUID(id) {
		return ctx.JSON(http.StatusBadRequest, util.PathVariableIsNotValid.ModifyApplicationName("user handler").ModifyErrorCode(4021))
	}

	res, errSrv := h.ticketService.TicketServiceDeleteById(id)
	if errSrv != nil {
		return ctx.JSON(errSrv.ErrorCode, util.NewError(errSrv.ApplicationName, errSrv.Operation, errSrv.Description, errSrv.ErrorCode, errSrv.StatusCode))
	}

	return ctx.JSON(http.StatusOK, res)
}
func (h *TicketHandler) TicketGetAll(ctx echo.Context) error {
	filter := util.Filter{}
	page, pageSize := util.ValidatePaginationFilters(ctx.QueryParam("page"), ctx.QueryParam("pageSize"), h.cfg.MaxPageLimit)
	filter.Page = page
	filter.PageSize = pageSize

	sortingField, sortingDirection := util.ValidateSortingFilters(entity.Category{}, ctx.QueryParam("sort"), ctx.QueryParam("sDirection"))
	filter.SortingField = sortingField
	filter.SortingDirection = sortingDirection

	filters := ctx.QueryParam("filters")
	filter.Filters = util.CreateFilter(entity.Ticket{}, filters)

	tickets, err := h.ticketService.TicketServiceGetAll(filter)
	if err != nil {
		return ctx.JSON(err.StatusCode, err)
	}

	if tickets == nil {
		return ctx.JSON(http.StatusNotFound, util.NotFound.ModifyApplicationName("user handler").ModifyErrorCode(5001))
	}
	return ctx.JSON(http.StatusOK, tickets)
}
