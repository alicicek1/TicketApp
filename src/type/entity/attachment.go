package entity

type Attachment struct {
	Id       string `json:"_id" bson:"_id,omitempty"`
	FileName string `json:"fileName" bson:"fileName,omitempty"`
	FilePath string `json:"filePath" bson:"filePath,omitempty"`
}

type AttachmentPostRequestModel struct {
	FileName string `json:"fileName" bson:"fileName,omitempty"`
	FilePath string `json:"filePath" bson:"filePath,omitempty"`
}

type AttachmentPostResponseModel struct {
	Id string `json:"_id" bson:"_id,omitempty"`
}

type AttachmentGetResponseModel struct {
	RowCount    int64  `json:"rowCount"`
	Attachments []User `json:"attachments"`
}
