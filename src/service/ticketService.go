package service

import (
	"github.com/google/uuid"
	"ticketApp/src/repository"
	"ticketApp/src/type/entity"
	"ticketApp/src/type/util"
	"time"
)

type TicketServiceType struct {
	TicketRepository repository.TicketRepository
}

type TicketService interface {
	TicketServiceInsert(user entity.Ticket) (*entity.TicketPostResponseModel, *util.Error)
	TicketServiceGetById(id string) (*entity.Ticket, *util.Error)
	TicketServiceDeleteById(id string) (util.DeleteResponseType, *util.Error)
	TicketServiceGetAll(filter util.Filter) (*entity.TicketGetReponseModel, *util.Error)
}

func NewTicketService(r repository.TicketRepository) TicketServiceType {
	return TicketServiceType{TicketRepository: r}
}

func (t TicketServiceType) TicketServiceInsert(ticket entity.Ticket) (*entity.TicketPostResponseModel, *util.Error) {
	if ticket.Id == "" {
		isSuccess, err := util.CheckTicketModel(ticket)
		if !isSuccess {
			return nil, err
		}
	}

	ticket.Id = uuid.New().String()
	ticket.CreatedAt = time.Now()
	ticket.UpdatedAt = time.Now()

	result, err := t.TicketRepository.TicketRepoInsert(ticket)

	return result, err
}
func (t TicketServiceType) TicketServiceGetById(id string) (*entity.Ticket, *util.Error) {
	result, err := t.TicketRepository.TicketRepoGetById(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (t TicketServiceType) TicketServiceDeleteById(id string) (util.DeleteResponseType, *util.Error) {
	result, err := t.TicketRepository.TicketRepoDeleteById(id)
	if err != nil || result.IsSuccess == false {
		return util.DeleteResponseType{IsSuccess: false}, err
	}
	return util.DeleteResponseType{IsSuccess: true}, nil
}
func (t TicketServiceType) TicketServiceGetAll(filter util.Filter) (*entity.TicketGetReponseModel, *util.Error) {
	result, err := t.TicketRepository.TicketRepositoryGetAll(filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
