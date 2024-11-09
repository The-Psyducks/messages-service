package messages

type RealTimeDatabaseInterface interface {
	SendMessage(senderId string, receiverId string, content string) (string, error)
	GetConversations() ([]string, error)
	SendNotificationToUserDevices(devicesTokens []string, title, body string) error
}
