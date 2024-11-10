package firebase_connector

type Interface interface {
	SendNotificationToUserDevices(devicesTokens []string, title, body string, data map[string]string) error
}
