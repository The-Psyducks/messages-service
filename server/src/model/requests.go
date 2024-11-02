package model

type MessageRequest struct {
	ReceiverId string `json:"receiver_id"`
	Content    string `json:"content"`
}

type NewDeviceRequest struct {
	DeviceId string `json:"device_id"`
}
