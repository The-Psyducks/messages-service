package repository

type DevicesDatabaseInterface interface {
	AddDevice(id string, token string) error
	GetDevicesTokens(id string) ([]string, error)
}
