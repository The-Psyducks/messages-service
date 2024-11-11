package firebase_connector

type MockFirebaseConnector struct {
}

func (m MockFirebaseConnector) SendNotificationToUserDevices(devicesTokens []string, title, body string, data map[string]string) error {
	return nil
}

func NewMockFirebaseConnector() Interface {
	return &MockFirebaseConnector{}
}
