package model

type MessageRequest struct {
	ReceiverId string `json:"receiver_id"`
	Content    string `json:"content"`
}

type NewDeviceRequest struct {
	DeviceId string `json:"device_id"`
}

type NotificationRequest struct {
	ReceiverId	string `json:"receiver_id"`
	Title		string `json:"title"`
	Body		string `json:"body"`
}