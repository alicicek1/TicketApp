package entity

import (
	"time"
)

type TicketStatus int16

const (
	CREATED  TicketStatus = 1
	ANSWERED              = 2
)

func (t TicketStatus) String() string {
	switch t {
	case CREATED:
		return "created"
	case ANSWERED:
		return "answered"
	}
	return "unknown"
}

type Ticket struct {
	Category       Category     `json:"category,omitempty" bson:"category,omitempty"`
	Attachments    []Attachment `json:"attachments,omitempty" bson:"attachments,omitempty"`
	Answers        []Answer     `json:"answers,omitempty" bson:"answers,omitempty"`
	Id             string       `json:"_id" bson:"_id,omitempty"`
	Subject        string       `json:"subject,omitempty" bson:"subject"`
	Body           string       `json:"body,omitempty" bson:"body"`
	CreatedBy      string       `json:"createdBy" bson:"createdBy,omitempty"`
	CreatedAt      time.Time    `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt      time.Time    `json:"updatedAt" bson:"updatedAt,omitempty"`
	LastAnsweredAt time.Time    `json:"lastAnsweredAt" bson:"lastAnsweredAt,omitempty"`
	Status         byte         `json:"status,omitempty" bson:"status,omitempty"`
}

type TicketPostRequestModel struct {
	Category       Category     `json:"category,omitempty"`
	Attachments    []Attachment `json:"attachments,omitempty"`
	Answers        []Answer     `json:"answers,omitempty"`
	Subject        string       `json:"subject,omitempty"`
	Body           string       `json:"body,omitempty"`
	CreatedBy      string       `json:"createdBy,omitempty"`
	LastAnsweredAt time.Time    `json:"lastAnsweredAt,omitempty"`
	Status         byte         `json:"status,omitempty,"`
}

type TicketPostResponseModel struct {
	Id string `json:"_id" bson:"_id,omitempty"`
}

type TicketGetReponseModel struct {
	RowCount int64    `json:"rowCount"`
	Tickets  []Ticket `json:"tickets"`
}
