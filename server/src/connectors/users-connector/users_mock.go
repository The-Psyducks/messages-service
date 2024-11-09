package users_connector

import (
	"errors"
)

type MockConnector struct {
	// Add fields as needed to simulate behavior
}

func NewMockConnector() Interface {
	return &MockConnector{}
}

func (m *MockConnector) CheckUserExists(userId, authHeader string) (bool, error) {

	switch userId {
	case "fakeUserId":
		return false, nil
	case "userId":
		return true, nil
	case "errorUserId":
		return true, errors.New("throwing error in mock")
	default:
		panic("id should be on of the following: fakeUserId, userId, errorUserId but it was: " + userId)
	}
}
