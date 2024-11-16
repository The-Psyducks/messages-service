package users_connector

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockConnector_GetUserNameAndImage(t *testing.T) {
	mock := NewMockConnector()

	username, image, err := mock.GetUserNameAndImage("someUserId", "someHeader")
	assert.NoError(t, err)
	assert.Equal(t, "someusername", username)
	assert.Equal(t, "someimage", image)
}

func TestMockConnector_CheckUserExists(t *testing.T) {
	mock := NewMockConnector()

	t.Run("UserExists", func(t *testing.T) {
		exists, err := mock.CheckUserExists("someUserId", "someHeader")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("UserDoesNotExist", func(t *testing.T) {
		exists, err := mock.CheckUserExists("fakeUserId", "someHeader")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("ErrorUserId", func(t *testing.T) {
		exists, err := mock.CheckUserExists("errorUserId", "someHeader")
		assert.Error(t, err)
		assert.True(t, exists)
		assert.Equal(t, errors.New("throwing error in mock"), err)
	})
}
