package repository

type MockDevicesDatabase struct {
}

func (m MockDevicesDatabase) AddDevice(_, _ string) error {
	//TODO implement me
	panic("implement me")
}

func (m MockDevicesDatabase) GetDevicesTokens(_ string) ([]string, error) {
	return []string{"token1", "token2"}, nil
}

func NewMockDevicesDatabase() DevicesDatabaseInterface {
	return &MockDevicesDatabase{}
}
