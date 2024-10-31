package model

type MessageDeliveredResponse struct {
	ChatReference string `json:"chat-reference"`
}

type GetMessagesResponse struct {
	ChatReferences []string `json:"chat-references"`
}
