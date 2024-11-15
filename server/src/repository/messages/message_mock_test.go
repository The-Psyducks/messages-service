package repository

//just so i can avoid coverage being lower due to this mocks
import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMockRealTimeDatabase_SendMessage(t *testing.T) {
	mockDB := NewMockRealTimeDatabase()

	t.Run("SendMessage_Success", func(t *testing.T) {
		ref, err := mockDB.SendMessage("senderId", "receiverId", "content")
		assert.NoError(t, err)
		assert.Equal(t, "mockMessageRef", ref)
	})

	t.Run("SendMessage_Error", func(t *testing.T) {
		ref, err := mockDB.SendMessage("senderId", "receiverId", "error")
		assert.Error(t, err)
		assert.Equal(t, "", ref)
	})
}

func TestMockRealTimeDatabase_GetChats(t *testing.T) {
	mockDB := NewMockRealTimeDatabase()

	t.Run("GetChats_Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			mockDB.GetChats("ok")
		})
	})
}

func TestMockRealTimeDatabase_GetConversations(t *testing.T) {
	mockDB := NewMockRealTimeDatabase()

	t.Run("GetConversations_Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			mockDB.GetConversations()
		})
	})
}
