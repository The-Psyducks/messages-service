package users_connector

import "errors"

type MockConnector struct {
	// Add fields as needed to simulate behavior
}

func (m *MockConnector) GetUserNameAndImage(id string, header string) (string, string, error) {
	return "some username", "some image", nil
}

func NewMockConnector() Interface {
	return &MockConnector{}
}

func (m *MockConnector) CheckUserExists(userId, authHeader string) (bool, error) {

	switch userId {
	case "fakeUserId":
		return false, nil
	default:
		return true, nil
	case "errorUserId":
		return true, errors.New("throwing error in mock")

	}

}
