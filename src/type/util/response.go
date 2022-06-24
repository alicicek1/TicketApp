package util

type DeleteResponseType struct {
	IsSuccess bool `json:"isSuccess"`
}

type PostResponseModel struct {
	Id string `json:"_id" bson:"_id,omitempty"`
}

type GetAllResponseType struct {
	RowCount int64 `json:"rowCount"`
	models   []any `json:"models"`
}
