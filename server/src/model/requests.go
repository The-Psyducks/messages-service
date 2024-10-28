package model

type MessageRequest struct {
	ReceiverId string `json:"receiver_id"`
	Content    string `json:"content"`
}
