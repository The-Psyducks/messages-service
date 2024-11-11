package model

type MessageDeliveredResponse struct {
	ChatReference string `json:"chat-reference"`
}

type GetMessagesResponse struct {
	ChatReferences []string `json:"chat-references"`
}

type ChatResponse struct {
	ChatReference string `json:"chat-reference"`
	UserName      string `json:"user-name"`
	UserImage     string `json:"user-image"`
	LastMessage   string `json:"last-message"`
	Date          string `json:"date"`
	toId          string `json:"to-id"`
}
