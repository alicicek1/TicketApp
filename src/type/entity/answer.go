package entity

import "time"

type Answer struct {
	Id        string    `json:"_id" bson:"_id,omitempty"`
	Body      string    `json:"body,omitempty" bson:"body,omitempty"`
	CreatedBy string    `json:"createdBy" bson:"createdBy,omitempty"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
}

type AnswerPostRequestModel struct {
	Body string `json:"body,omitempty" bson:"body,omitempty"`
}

type AnswerPostResponseModel struct {
	Id string `json:"_id" bson:"_id,omitempty"`
}

type AnswerGetResponseModel struct {
	RowCount int64  `json:"rowCount"`
	Answers  []User `json:"answers"`
}
